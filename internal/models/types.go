// The `models` package defines several data structures that represent responses from Ethereum-related APIs.
// These structures are used to unmarshal JSON responses into Go objects for further processing in the application.

package models

// BeaconBlockResponse represents the response structure for a beacon block request.
// It contains nested structs to capture the version and execution payload details of the block.
type BeaconBlockResponse struct {
	Version string `json:"version"` // The version of the beacon block.
	Data    struct {
		Message struct {
			Body struct {
				ExecutionPayload struct {
					BlockNumber   string `json:"block_number"`     // The block number in the execution payload.
					FeeRecipient  string `json:"fee_recipient"`    // The address that receives the transaction fees.
					ExtraData     string `json:"extra_data"`       // Additional data included in the block.
					BaseFeePerGas string `json:"base_fee_per_gas"` // The base fee per gas unit for the block.
					GasUsed       string `json:"gas_used"`         // The total gas used by transactions in the block.
				} `json:"execution_payload"`
			} `json:"body"`
		} `json:"message"`
	} `json:"data"`
}

// BeaconHeadersResponse represents the response structure for beacon headers.
// It includes a list of headers, each containing a message with a slot identifier.
type BeaconHeadersResponse struct {
	Data []struct {
		Header struct {
			Message struct {
				Slot string `json:"slot"` // The slot number associated with the beacon header.
			} `json:"message"`
		} `json:"header"`
	} `json:"data"`
}

// ExecutionBlockTx represents a transaction within an execution block.
// It includes various fields such as block hash, gas details, and transaction identifiers.
type ExecutionBlockTx struct {
	BlockHash        string `json:"blockHash"`        // The hash of the block containing the transaction.
	BlockNumber      string `json:"blockNumber"`      // The block number containing the transaction.
	From             string `json:"from"`             // The address that initiated the transaction.
	Gas              string `json:"gas"`              // The gas limit provided by the sender.
	GasPrice         string `json:"gasPrice"`         // The price per gas unit offered by the sender.
	Hash             string `json:"hash"`             // The hash of the transaction.
	Input            string `json:"input"`            // The input data for the transaction.
	Nonce            string `json:"nonce"`            // The number of transactions sent from the sender's address.
	To               string `json:"to"`               // The address of the recipient.
	TransactionIndex string `json:"transactionIndex"` // The index of the transaction within the block.
	Value            string `json:"value"`            // The amount of Ether transferred.
	Type             string `json:"type"`             // The type of transaction.
}

// ExecutionBlockFullResponse represents the full response for an execution block request.
// It includes the block number, base fee, extra data, and a list of transactions.
type ExecutionBlockFullResponse struct {
	Result struct {
		Number        string             `json:"number"`        // The block number.
		BaseFeePerGas string             `json:"baseFeePerGas"` // The base fee per gas unit for the block.
		ExtraData     string             `json:"extraData"`     // Additional data included in the block.
		Transactions  []ExecutionBlockTx `json:"transactions"`  // A list of transactions in the block.
	} `json:"result"`
}

// SyncCommitteeResponse represents the response from the sync_committees endpoint.
// It includes flags for execution optimism and finalization, along with a list of validator addresses.
type SyncCommitteeResponse struct {
	ExecutionOptimistic bool `json:"execution_optimistic"` // Indicates if the execution is optimistic.
	Finalized           bool `json:"finalized"`            // Indicates if the data is finalized.
	Data                struct {
		Validators []string `json:"validators"` // A list of validator addresses in the sync committee.
	} `json:"data"`
}
