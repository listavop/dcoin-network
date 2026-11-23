
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
	Token      *ledger.TestToken
}

func NewServer(
	bc *core.Blockchain,
	l *ledger.Ledger,
	token *ledger.TestToken,
) *Server {
	return &Server{
		Blockchain: bc,
		Ledger:     l,
		Token:      token,
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status": "ok",
			"network": "DCOIN Testnet",
		})
	})

	mux.HandleFunc("/api/v1/network", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"name":         "DCOIN Testnet",
			"symbol":       "DCO",
			"chain_id":     "dcoin-testnet-1",
			"decimals":     ledger.Decimals,
			"total_supply": ledger.TotalSupply,
			"block_height": s.Blockchain.LatestBlock().Index,
			"mining":       false,
		})
	})

	mux.HandleFunc("/api/v1/blocks/latest", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, s.Blockchain.LatestBlock())
	})

	mux.HandleFunc("/api/v1/balance", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "address is required",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"address": address,
			"balance": s.Ledger.Balance(address),
			"symbol":  "DCO",
			"decimals": ledger.Decimals,
		})
	})

	mux.HandleFunc("/api/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, s.Ledger.TransactionHistory())
	})

	mux.HandleFunc("/api/v1/token", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"name":        ledger.TokenName,
			"symbol":      ledger.TokenSymbol,
			"decimals":    ledger.TokenDecimals,
			"max_supply":  ledger.TokenMaxSupply,
			"total_supply": s.Token.Supply,
		})
	})

	mux.HandleFunc("/api/v1/token/balance", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "address is required",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"address": address,
			"balance": s.Token.Balance(address),
			"symbol":  ledger.TokenSymbol,
		})
	})

	mux.HandleFunc("/api/v1/token/faucet", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
				"error": "POST required",
			})
			return
		}

		var request struct {
			Address string `json:"address"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid JSON",
			})
			return
		}

		if err := s.Token.Faucet(request.Address); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "success",
			"address": request.Address,
			"amount":  ledger.FaucetAmount,
			"symbol":  ledger.TokenSymbol,
		})
	})

	return mux
}
