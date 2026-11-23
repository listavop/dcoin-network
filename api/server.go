
package api

import (
	"encoding/json"
	"net/http"

	"github.com/listavop/dcoin-network/core"
	"github.com/listavop/dcoin-network/ledger"
)

type Server struct {
	Blockchain *core.Blockchain
	Ledger     *ledger.Ledger
}

func NewServer(
	blockchain *core.Blockchain,
	accountLedger *ledger.Ledger,
) *Server {
	return &Server{
		Blockchain: blockchain,
		Ledger:     accountLedger,
	}
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/api/v1/network", s.network)
	mux.HandleFunc("/api/v1/blocks/latest", s.latestBlock)
	mux.HandleFunc("/api/v1/balance", s.balance)
	mux.HandleFunc("/api/v1/transactions", s.transactions)

	return mux
}

func (s *Server) health(
	w http.ResponseWriter,
	r *http.Request,
) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "online",
		"network": "DCOIN Testnet",
	})
}

func (s *Server) network(
	w http.ResponseWriter,
	r *http.Request,
) {
	if !s.Blockchain.Validate() {
		writeJSON(w, http.StatusInternalServerError,
			map[string]string{"error": "blockchain validation failed"})
		return
	}

	latest := s.Blockchain.LatestBlock()

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"name":         "DCOIN Testnet",
		"symbol":       "DCO",
		"chain_id":     "dcoin-testnet-1",
		"decimals":     8,
		"total_supply": ledger.TotalSupply,
		"block_height": latest.Index,
		"latest_hash":  latest.Hash,
		"mining":       false,
	})
}

func (s *Server) latestBlock(
	w http.ResponseWriter,
	r *http.Request,
) {
	block := s.Blockchain.LatestBlock()
	writeJSON(w, http.StatusOK, block)
}

func (s *Server) balance(
	w http.ResponseWriter,
	r *http.Request,
) {
	address := r.URL.Query().Get("address")

	if address == "" {
		writeJSON(w, http.StatusBadRequest,
			map[string]string{"error": "address required"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"address": address,
		"balance": s.Ledger.Balance(address),
		"asset":   "DCO",
		"decimals": ledger.Decimals,
	})
}

func (s *Server) transactions(
	w http.ResponseWriter,
	r *http.Request,
) {
	history := s.Ledger.TransactionHistory()

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count":        len(history),
		"transactions": history,
	})
}
