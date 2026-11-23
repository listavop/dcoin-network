
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
	log.Println("Starting DCOIN Testnet...")

	dataPath := os.Getenv("DCOIN_DATA")
	if dataPath == "" {
		dataPath = "data/dcoin.json"
	}

	var (
		blockchain    *core.Blockchain
		accountLedger *ledger.Ledger
		err           error
	)

	blockchain, accountLedger, err = storage.Load(dataPath)

	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal("Storage loading failed:", err)
		}

		log.Println("Creating new DCOIN blockchain...")

		blockchain = core.NewBlockchain()

		accountLedger, err = ledger.NewLedger("DCOIN_GENESIS")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Existing blockchain loaded")
	}

	validator, err := core.NewValidator()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Validator ID:", validator.ID)

	server := api.NewServer(blockchain, accountLedger)

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
		log.Println("DCOIN RPC listening on port", port)
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
		log.Println("Stopping DCOIN Node...")
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Println("Server error:", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Println("Shutdown error:", err)
	}

	if err := storage.Save(
		dataPath,
		blockchain,
		accountLedger,
	); err != nil {
		log.Println("Storage save failed:", err)
	} else {
		log.Println("DCOIN blockchain saved successfully")
	}

	log.Println("DCOIN Node stopped")
}
