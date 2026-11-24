
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/listavop/dcoin-network/core"
	"github.com/listavop/dcoin-network/ledger"
)

func newTestServer(t *testing.T) http.Handler {
	t.Helper()

	bc := core.NewBlockchain()

	accountLedger, err := ledger.NewLedger("DCOIN_GENESIS")
	if err != nil {
		t.Fatal(err)
	}

	token := ledger.NewTestToken()

	server := NewServer(bc, accountLedger, token)
	return server.Handler()
}

func TestHealthEndpoint(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestNetworkEndpoint(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/network", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestTokenEndpoint(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/token", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestTokenFaucet(t *testing.T) {
	handler := newTestServer(t)

	body := map[string]string{
		"address": "DCOIN_TEST_WALLET",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/token/faucet",
		bytes.NewReader(data),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	tokenBalanceReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/token/balance?address=DCOIN_TEST_WALLET",
		nil,
	)

	tokenBalanceRec := httptest.NewRecorder()
	handler.ServeHTTP(tokenBalanceRec, tokenBalanceReq)

	if tokenBalanceRec.Code != http.StatusOK {
		t.Fatalf("expected balance status 200, got %d", tokenBalanceRec.Code)
	}
}

func TestFaucetCannotBeClaimedTwice(t *testing.T) {
	handler := newTestServer(t)

	for i := 0; i < 2; i++ {
		body := []byte(`{"address":"DCOIN_TEST_WALLET"}`)

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/token/faucet",
			bytes.NewReader(body),
		)

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if i == 0 && rec.Code != http.StatusOK {
			t.Fatalf("first claim failed: %d", rec.Code)
		}

		if i == 1 && rec.Code != http.StatusBadRequest {
			t.Fatalf("second claim should fail, got %d", rec.Code)
		}
	}
}
