
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/listavop/dcoin-network/api"
	"github.com/listavop/dcoin-network/core"
	"github.com/listavop/dcoin-network/ledger"
)

func main() {
	log.Println("Starting DCOIN Testnet Node...")

	// Initialize blockchain
	blockchain := core.NewBlockchain()

	// Initialize genesis ledger
	accountLedger, err := ledger.NewLedger("DCOIN_GENESIS")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize validator identity
	validator, err := core.NewValidator()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DCOIN Validator initialized")
	log.Println("Validator ID:", validator.ID)

	// Initialize API
	server := api.NewServer(blockchain, accountLedger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: server.Handler(),
	}

	log.Println("Network: DCOIN Testnet")
	log.Println("Chain ID: dcoin-testnet-1")
	log.Println("Genesis Block:", blockchain.LatestBlock().Hash)
	log.Println("RPC Server port:", port)

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal("DCOIN server error:", err)
	}
}
