package ledger

import (
	"errors"
	"math"
	"sync"
)

const (
	TokenName       = "DCOIN Test Token"
	TokenSymbol     = "TDCO"
	TokenDecimals   = 8
	TokenMaxSupply  = uint64(1_000_000_000) * 100_000_000
	FaucetAmount    = uint64(1_000) * 100_000_000
)

type TestToken struct {
	mu       sync.RWMutex
	Balances map[string]uint64 `json:"balances"`
	Claimed  map[string]bool   `json:"claimed"`
	Supply   uint64            `json:"supply"`
}

func NewTestToken() *TestToken {
	return &TestToken{
		Balances: make(map[string]uint64),
		Claimed:  make(map[string]bool),
		Supply:   0,
	}
}

func (t *TestToken) Faucet(address string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if address == "" {
		return errors.New("invalid address")
	}

	if t.Claimed[address] {
		return errors.New("faucet already claimed")
	}

	if t.Supply > TokenMaxSupply-FaucetAmount {
		return errors.New("token faucet supply exhausted")
	}

	balance := t.Balances[address]

	if math.MaxUint64-balance < FaucetAmount {
		return errors.New("balance overflow")
	}

	t.Balances[address] = balance + FaucetAmount
	t.Supply += FaucetAmount
	t.Claimed[address] = true

	return nil
}

func (t *TestToken) Balance(address string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.Balances[address]
}

func (t *TestToken) Transfer(
	from string,
	to string,
	amount uint64,
) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if from == "" || to == "" || from == to {
		return errors.New("invalid address")
	}

	if amount == 0 {
		return errors.New("invalid amount")
	}

	if t.Balances[from] < amount {
		return errors.New("insufficient TDCO balance")
	}

	if math.MaxUint64-t.Balances[to] < amount {
		return errors.New("balance overflow")
	}

	t.Balances[from] -= amount
	t.Balances[to] += amount

type TestTokenState struct {
	Balances map[string]uint64 `json:"balances"`
	Claimed  map[string]bool   `json:"claimed"`
	Supply   uint64            `json:"supply"`
}

func (t *TestToken) Snapshot() TestTokenState {
	t.mu.RLock()
	defer t.mu.RUnlock()

	balances := make(map[string]uint64, len(t.Balances))
	for address, amount := range t.Balances {
		balances[address] = amount
	}

	claimed := make(map[string]bool, len(t.Claimed))
	for address, value := range t.Claimed {
		claimed[address] = value
	}

	return TestTokenState{
		Balances: balances,
		Claimed:  claimed,
		Supply:   t.Supply,
	}
}

func RestoreTestToken(state TestTokenState) (*TestToken, error) {
	if state.Supply > TokenMaxSupply {
		return nil, errors.New("token supply exceeds maximum")
	}

	var total uint64
	for _, amount := range state.Balances {
		if math.MaxUint64-total < amount {
			return nil, errors.New("token balance overflow")
		}
		total += amount
	}

	if total != state.Supply {
		return nil, errors.New("token supply does not match balances")
	}

	token := NewTestToken()
	token.Balances = state.Balances
	token.Claimed = state.Claimed
	token.Supply = state.Supply

	return token, nil
}
	
func (t *TestToken) Transfer(from, to string, amount uint64) error {
	// Isi fungsi Transfer yang sudah ada tetap dipertahankan.
	return nil
}

type TestTokenState struct {
	Balances map[string]uint64 `json:"balances"`
	Claimed  map[string]bool   `json:"claimed"`
	Supply   uint64            `json:"supply"`
}
	

