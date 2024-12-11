// This Go program is the entry point for a web server that provides Ethereum rewards-related services.
// It uses the Gin web framework to handle HTTP requests and the godotenv package to load environment variables from a .env file.

package main

import (
	"eth-rewards-api/internal/handlers"
	"eth-rewards-api/internal/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // For loading .env file
)

func main() {
	// Attempt to load environment variables from a .env file.
	// If the file is not found or fails to load, log a message but continue execution.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load.")
	}

	// Retrieve the QUICKNODE_ENDPOINT environment variable, which is expected to contain the endpoint URL.
	// If the variable is not set, log a fatal error and terminate the program.
	endpoint := os.Getenv("QUICKNODE_ENDPOINT")
	if endpoint == "" {
		log.Fatal("QUICKNODE_ENDPOINT environment variable not set.")
	}

	// Initialize services for consensus and execution layers using the endpoint.
	consensusService := services.NewConsensusService(endpoint)
	executionService := services.NewExecutionService(endpoint)

	// Create a new Gin router instance.
	r := gin.Default()

	// Create a new BlockRewardHandler with the initialized services.
	blockRewardHandler := handlers.NewBlockRewardHandler(consensusService, executionService)

	// Define an HTTP GET endpoint for retrieving block rewards by slot.
	r.GET("/blockreward/:slot", blockRewardHandler.GetBlockReward)

	// Define an HTTP GET endpoint for retrieving sync committee duties by slot.
	r.GET("/syncduties/:slot", blockRewardHandler.GetSyncDuties)

	// Start the Gin server on port 8080.
	// If the server fails to start, log a fatal error and terminate the program.
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
