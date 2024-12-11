// The `services` package provides functionality to interact with Ethereum consensus layer APIs.
// It includes a `ConsensusService` struct that handles HTTP requests to fetch data related to beacon chain slots and sync committees.

package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"eth-rewards-api/internal/models"
)

// SLOTS_PER_EPOCH is a constant that defines the number of slots in a single epoch on the Ethereum mainnet.
const SLOTS_PER_EPOCH = 32

// ConsensusService is a struct that holds the endpoint URL and an HTTP client for making requests.
type ConsensusService struct {
	endpoint string
	client   *http.Client
}

// NewConsensusService initializes a new instance of ConsensusService with a specified endpoint and a default HTTP client.
func NewConsensusService(endpoint string) *ConsensusService {
	return &ConsensusService{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 10 * time.Second, // Sets a timeout for HTTP requests.
		},
	}
}

// GetHeadSlot retrieves the current head slot number from the beacon chain headers endpoint.
// It returns the slot number as a uint64 and an error if any issues occur during the request or data parsing.
func (c *ConsensusService) GetHeadSlot() (uint64, error) {
	url := fmt.Sprintf("%s/eth/v1/beacon/headers", c.endpoint)
	resp, err := c.client.Get(url)
	if err != nil {
		return 0, err // Return an error if the HTTP request fails.
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode) // Handle non-200 HTTP responses.
	}

	var headersResp models.BeaconHeadersResponse
	if err := json.NewDecoder(resp.Body).Decode(&headersResp); err != nil {
		return 0, err // Return an error if JSON decoding fails.
	}
	if len(headersResp.Data) == 0 {
		return 0, errors.New("no header data returned") // Handle empty data response.
	}
	headSlotStr := headersResp.Data[0].Header.Message.Slot
	headSlot, err := strconv.ParseUint(headSlotStr, 10, 64)
	if err != nil {
		return 0, err // Return an error if slot conversion fails.
	}
	return headSlot, nil // Return the head slot number.
}

// GetBeaconBlockBySlot fetches the beacon block for a given slot number.
// It returns a pointer to a BeaconBlockResponse and an error if any issues occur during the request or data parsing.
func (c *ConsensusService) GetBeaconBlockBySlot(slot uint64) (*models.BeaconBlockResponse, error) {
	url := fmt.Sprintf("%s/eth/v2/beacon/blocks/%d", c.endpoint, slot)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err // Return an error if the HTTP request fails.
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("block not found") // Handle 404 response.
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode) // Handle non-200 HTTP responses.
	}

	var blockResp models.BeaconBlockResponse
	if err := json.NewDecoder(resp.Body).Decode(&blockResp); err != nil {
		return nil, err // Return an error if JSON decoding fails.
	}
	return &blockResp, nil // Return the beacon block response.
}

// GetSyncCommitteeDuties retrieves the sync committee validators for a specified slot.
// It calculates the epoch and constructs the state_id to fetch the relevant data.
// Returns a slice of validator addresses and an error if any issues occur during the request or data parsing.
func (c *ConsensusService) GetSyncCommitteeDuties(slot uint64) ([]string, error) {
	epoch := slot / SLOTS_PER_EPOCH
	state_id := epoch * SLOTS_PER_EPOCH // Calculate the first slot of the epoch.
	url := fmt.Sprintf("%s/eth/v1/beacon/states/%d/sync_committees?epoch=%d", c.endpoint, state_id, epoch)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err // Return an error if the HTTP request fails.
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("sync committee duties not found for this slot") // Handle 404 response.
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from sync duties endpoint", resp.StatusCode) // Handle non-200 HTTP responses.
	}

	var scResp models.SyncCommitteeResponse
	if err := json.NewDecoder(resp.Body).Decode(&scResp); err != nil {
		return nil, err // Return an error if JSON decoding fails.
	}

	return scResp.Data.Validators, nil // Return the list of validator addresses.
}
