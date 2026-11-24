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
	Blockchain core.Blockchain      `json:"blockchain"`
	Ledger     ledger.Ledger        `json:"ledger"`
	Token      ledger.TestTokenState `json:"token"`
}

func Save(path string, bc *core.Blockchain, l *ledger.Ledger) error {
	return SaveWithToken(path, bc, l, ledger.NewTestToken())
}

func SaveWithToken(
	path string,
	bc *core.Blockchain,
	l *ledger.Ledger,
	token *ledger.TestToken,
) error {
	if bc == nil || l == nil || token == nil {
		return errors.New("storage: nil state")
	}

	if !bc.Validate() {
		return errors.New("storage: invalid blockchain")
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
		Token: token.Snapshot(),
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	tempPath := path + ".tmp"

	if err := os.WriteFile(tempPath, data, 0600); err != nil {
		return err
	}

	if err := os.Rename(tempPath, path); err != nil {
		_ = os.Remove(tempPath)
		return err
	}

	return nil
}

func Load(path string) (*core.Blockchain, *ledger.Ledger, error) {
	bc, accountLedger, _, err := LoadWithToken(path)
	return bc, accountLedger, err
}

func LoadWithToken(
	path string,
) (*core.Blockchain, *ledger.Ledger, *ledger.TestToken, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, nil, nil, err
	}

	bc := &core.Blockchain{
		Blocks: state.Blockchain.Blocks,
	}

	if !bc.Validate() {
		return nil, nil, nil, errors.New("storage: invalid blockchain")
	}

	if state.Ledger.TotalSupply != ledger.TotalSupply {
		return nil, nil, nil, errors.New("storage: incorrect DCOIN total supply")
	}

	var total uint64

	for _, amount := range state.Ledger.Balances {
		if math.MaxUint64-total < amount {
			return nil, nil, nil, errors.New("storage: balance overflow")
		}
		total += amount
	}

	if total != ledger.TotalSupply {
		return nil, nil, nil, errors.New("storage: DCOIN balances do not match supply")
	}

	accountLedger := &ledger.Ledger{
		Balances:     state.Ledger.Balances,
		Transactions: state.Ledger.Transactions,
		TotalSupply:  state.Ledger.TotalSupply,
	}

	token, err := ledger.RestoreTestToken(state.Token)
	if err != nil {
		return nil, nil, nil, err
	}

	return bc, accountLedger, token, nil
}
