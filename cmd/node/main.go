
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/listavop/dcoin-network/api"
	"github.com/listavop/dcoin-network/core"
	"github.com/listavop/dcoin-network/ledger"
	"github.com/listavop/dcoin-network/storage"
)

func main() {
	dataPath := os.Getenv("DCOIN_DATA")
	if dataPath == "" {
		dataPath = "data/dcoin.json"
	}

	var blockchain *core.Blockchain
	var accountLedger *ledger.Ledger
	var testToken *ledger.TestToken
	var err error

	blockchain, accountLedger, testToken, err = storage.LoadWithToken(dataPath)

	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatalf("Failed to load DCOIN state: %v", err)
		}

		log.Println("Creating new DCOIN testnet state")

		blockchain = core.NewBlockchain()

		accountLedger, err = ledger.NewLedger("DCOIN_GENESIS")
		if err != nil {
			log.Fatalf("Failed to create ledger: %v", err)
		}

		testToken = ledger.NewTestToken()
	}

	validator, err := core.NewValidator()
	if err != nil {
		log.Fatalf("Failed to create validator: %v", err)
	}

	log.Printf("Validator ID: %s", validator.ID)

	server := api.NewServer(blockchain, accountLedger, testToken)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           server.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("DCOIN Testnet API listening on port %s", port)
		serverErrors <- httpServer.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")

	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}

	if err := storage.SaveWithToken(
		dataPath,
		blockchain,
		accountLedger,
		testToken,
	); err != nil {
		log.Printf("Failed to save DCOIN state: %v", err)
	} else {
		log.Println("DCOIN state saved successfully")
	}

	log.Println("DCOIN node stopped")
}

