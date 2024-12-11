// This package defines handlers for processing HTTP requests related to Ethereum block rewards and sync committee duties.
package handlers

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"eth-rewards-api/internal/services"

	"github.com/gin-gonic/gin"
)

// BlockRewardHandler is a struct that holds references to the consensus and execution services.
type BlockRewardHandler struct {
	consensusService *services.ConsensusService
	executionService *services.ExecutionService
}

// NewBlockRewardHandler initializes a new BlockRewardHandler with the provided services.
func NewBlockRewardHandler(cs *services.ConsensusService, es *services.ExecutionService) *BlockRewardHandler {
	return &BlockRewardHandler{
		consensusService: cs,
		executionService: es,
	}
}

// GetBlockReward handles HTTP requests to retrieve the block reward for a given slot.
func (h *BlockRewardHandler) GetBlockReward(c *gin.Context) {
	// Parse the slot parameter from the request URL.
	slotParam := c.Param("slot")
	slot, err := strconv.ParseUint(slotParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot parameter"})
		return
	}

	// Ensure the requested slot is not in the future by comparing it with the current head slot.
	headSlot, err := h.consensusService.GetHeadSlot()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch head slot"})
		return
	}
	if slot > headSlot {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requested slot is in the future"})
		return
	}

	// Retrieve the beacon block for the specified slot.
	beaconBlock, err := h.consensusService.GetBeaconBlockBySlot(slot)
	if err != nil {
		if err.Error() == "block not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "slot not found/missed"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get beacon block"})
		return
	}

	// Extract the block number from the beacon block's execution payload.
	blockNumberDecimal := beaconBlock.Data.Message.Body.ExecutionPayload.BlockNumber
	if blockNumberDecimal == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "no execution payload for this slot"})
		return
	}

	// Convert the block number to hexadecimal format.
	blockNumberInt, err := strconv.ParseUint(blockNumberDecimal, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid block number format"})
		return
	}
	blockNumberHex := fmt.Sprintf("0x%x", blockNumberInt)

	// Retrieve the execution block using the block number in hexadecimal format.
	execBlock, err := h.executionService.GetExecutionBlockByNumber(blockNumberHex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get execution block"})
		return
	}

	// Calculate the total reward by iterating over each transaction in the execution block.
	baseFee, err := hexToBigInt(execBlock.Result.BaseFeePerGas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid base fee"})
		return
	}

	totalReward := big.NewInt(0)
	for _, tx := range execBlock.Result.Transactions {
		gasPrice, err := hexToBigInt(tx.GasPrice)
		if err != nil {
			continue
		}
		gas, err := hexToBigInt(tx.Gas)
		if err != nil {
			continue
		}

		// Calculate the transaction reward if the gas price is greater than the base fee.
		if gasPrice.Cmp(baseFee) > 0 {
			priorityFee := big.NewInt(0).Sub(gasPrice, baseFee)
			txReward := big.NewInt(0).Mul(priorityFee, gas)
			totalReward.Add(totalReward, txReward)
		}
	}

	// Convert the total reward from wei to gwei.
	divider := big.NewInt(1_000_000_000)
	rewardInGwei := big.NewInt(0).Div(totalReward, divider)

	// Determine the status based on the length of the extra data in the execution block.
	status := "vanilla"
	if len(execBlock.Result.ExtraData) > 20 {
		status = "relay"
	}

	// Respond with the calculated reward and status.
	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"reward": rewardInGwei.String(),
	})
}

// GetSyncDuties handles HTTP requests to retrieve sync committee duties for a given slot.
func (h *BlockRewardHandler) GetSyncDuties(c *gin.Context) {
	// Parse the slot parameter from the request URL.
	slotParam := c.Param("slot")
	slot, err := strconv.ParseUint(slotParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot parameter"})
		return
	}

	// Ensure the requested slot is not too far in the future by comparing it with the current head slot.
	headSlot, err := h.consensusService.GetHeadSlot()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch head slot"})
		return
	}
	if slot > headSlot {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requested slot is too far in the future"})
		return
	}

	// Retrieve the sync committee duties for the specified slot.
	validators, err := h.consensusService.GetSyncCommitteeDuties(slot)
	if err != nil {
		if err.Error() == "sync committee duties not found for this slot" {
			c.JSON(http.StatusNotFound, gin.H{"error": "sync committee duties not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sync committee duties"})
		return
	}

	// Respond with the list of validators in the sync committee.
	c.JSON(http.StatusOK, gin.H{
		"validators": validators,
	})
}

// hexToBigInt converts a hexadecimal string to a big.Int.
func hexToBigInt(hexStr string) (*big.Int, error) {
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		i := new(big.Int)
		_, ok := i.SetString(hexStr[2:], 16)
		if !ok {
			return nil, fmt.Errorf("failed to parse hex string")
		}
		return i, nil
	}
	return nil, fmt.Errorf("invalid hex format")
}
