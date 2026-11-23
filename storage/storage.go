
package storage

import (
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"

	"github.com/listavop/dcoin-network/core"
	"github.com/listavop/dcoin-network/ledger"
)

type State struct {
	Blockchain core.Blockchain `json:"blockchain"`
	Ledger     ledger.Ledger   `json:"ledger"`
}

func Save(
	path string,
	bc *core.Blockchain,
	l *ledger.Ledger,
) error {
	if bc == nil || l == nil {
		return errors.New("blockchain or ledger is nil")
	}

	if !bc.Validate() {
		return errors.New("invalid blockchain")
	}

	state := State{
		Blockchain: core.Blockchain{
			Blocks: append([]core.Block(nil), bc.Blocks...),
		},
		Ledger: ledger.Ledger{
			Balances:     l.Balances,
			Transactions: l.TransactionHistory(),
			TotalSupply:  l.TotalSupply,
		},
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	temp := path + ".tmp"

	if err := os.WriteFile(temp, data, 0600); err != nil {
		return err
	}

	if err := os.Rename(temp, path); err != nil {
		os.Remove(temp)
		return err
	}

	return nil
}

func Load(
	path string,
) (*core.Blockchain, *ledger.Ledger, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	var state State

	if err := json.Unmarshal(data, &state); err != nil {
		return nil, nil, err
	}

	bc := &core.Blockchain{
		Blocks: state.Blockchain.Blocks,
	}

	if !bc.Validate() {
		return nil, nil, errors.New("blockchain integrity check failed")
	}

	if state.Ledger.TotalSupply != ledger.TotalSupply {
		return nil, nil, errors.New("invalid total supply")
	}

	var total uint64

	for _, balance := range state.Ledger.Balances {
		if math.MaxUint64-total < balance {
			return nil, nil, errors.New("balance overflow")
		}
		total += balance
	}

	if total != ledger.TotalSupply {
		return nil, nil, errors.New("ledger supply mismatch")
	}

	accountLedger := &ledger.Ledger{
		Balances:     state.Ledger.Balances,
		Transactions: state.Ledger.Transactions,
		TotalSupply:  state.Ledger.TotalSupply,
	}

	return bc, accountLedger, nil
}
