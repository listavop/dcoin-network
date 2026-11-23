
package ledger

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	Decimals   uint64 = 8
	TotalSupply uint64 = 7_000_000 * 100_000_000
)

type Transaction struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    uint64 `json:"amount"`
	Timestamp int64  `json:"timestamp"`
	Asset     string `json:"asset"`
}

type Ledger struct {
	mu           sync.RWMutex
	Balances     map[string]uint64 `json:"balances"`
	Transactions []Transaction     `json:"transactions"`
	TotalSupply  uint64            `json:"total_supply"`
}

func NewLedger(genesisAddress string) (*Ledger, error) {
	if genesisAddress == "" {
		return nil, errors.New("genesis address required")
	}

	balances := make(map[string]uint64)
	balances[genesisAddress] = TotalSupply

	return &Ledger{
		Balances:     balances,
		Transactions: []Transaction{},
		TotalSupply:  TotalSupply,
	}, nil
}

func (l *Ledger) Balance(address string) uint64 {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.Balances[address]
}

func (l *Ledger) Transfer(
	from string,
	to string,
	amount uint64,
) (Transaction, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if from == "" || to == "" {
		return Transaction{}, errors.New("invalid address")
	}

	if from == to {
		return Transaction{}, errors.New("self transfer not allowed")
	}

	if amount == 0 {
		return Transaction{}, errors.New("amount must be greater than zero")
	}

	senderBalance := l.Balances[from]
	receiverBalance := l.Balances[to]

	if senderBalance < amount {
		return Transaction{}, errors.New("insufficient balance")
	}

	if math.MaxUint64-receiverBalance < amount {
		return Transaction{}, errors.New("balance overflow")
	}

	timestamp := time.Now().UnixNano()

	record := fmt.Sprintf(
		"%s|%s|%d|%d|%d",
		from,
		to,
		amount,
		timestamp,
		len(l.Transactions),
	)

	hash := sha256.Sum256([]byte(record))

	tx := Transaction{
		ID:        hex.EncodeToString(hash[:]),
		From:      from,
		To:        to,
		Amount:    amount,
		Timestamp: timestamp,
		Asset:     "DCO",
	}

	l.Balances[from] = senderBalance - amount
	l.Balances[to] = receiverBalance + amount

	l.Transactions = append(l.Transactions, tx)

	return tx, nil
}

func (l *Ledger) TransactionHistory() []Transaction {
	l.mu.RLock()
	defer l.mu.RUnlock()

	history := make([]Transaction, len(l.Transactions))
	copy(history, l.Transactions)

	return history
}
