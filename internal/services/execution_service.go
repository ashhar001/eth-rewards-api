// The `services` package provides functionality to interact with Ethereum execution layer APIs.
// It includes an `ExecutionService` struct that handles JSON-RPC requests to fetch execution block data.

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"eth-rewards-api/internal/models"
)

// ExecutionService is a struct that holds the endpoint URL and an HTTP client for making requests.
type ExecutionService struct {
	endpoint string
	client   *http.Client
}

// NewExecutionService initializes a new instance of ExecutionService with a specified endpoint and a default HTTP client.
func NewExecutionService(endpoint string) *ExecutionService {
	return &ExecutionService{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 10 * time.Second, // Sets a timeout for HTTP requests.
		},
	}
}

// JSONRPCRequest represents the structure of a JSON-RPC request.
// It includes the JSON-RPC version, method name, parameters, and an identifier.
type JSONRPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

// GetExecutionBlockByNumber sends a JSON-RPC request to retrieve an execution block by its number in hexadecimal format.
// It returns a pointer to an ExecutionBlockFullResponse and an error if any issues occur during the request or data parsing.
func (e *ExecutionService) GetExecutionBlockByNumber(blockNumberHex string) (*models.ExecutionBlockFullResponse, error) {
	// Create a JSON-RPC request body with the method "eth_getBlockByNumber" and the block number as a parameter.
	reqBody := JSONRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{blockNumberHex, true},
		Id:      1,
	}
	// Marshal the request body into JSON format.
	b, _ := json.Marshal(reqBody)
	// Send a POST request to the execution endpoint with the JSON-RPC request body.
	resp, err := e.client.Post(e.endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err // Return an error if the HTTP request fails.
	}
	defer resp.Body.Close()

	// Check if the response status code is not 200 OK.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode) // Handle non-200 HTTP responses.
	}

	// Decode the JSON response body into an ExecutionBlockFullResponse struct.
	var blockResp models.ExecutionBlockFullResponse
	if err := json.NewDecoder(resp.Body).Decode(&blockResp); err != nil {
		return nil, err // Return an error if JSON decoding fails.
	}
	// Check if the block number in the response is empty, indicating the block was not found.
	if blockResp.Result.Number == "" {
		return nil, fmt.Errorf("block not found on execution layer") // Handle block not found scenario.
	}
	return &blockResp, nil // Return the execution block response.
}
