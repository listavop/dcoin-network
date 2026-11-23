
package core

import (
	"errors"
	"sync"
)

type Blockchain struct {
	mu     sync.RWMutex
	Blocks []Block `json:"blocks"`
}

func NewBlockchain() *Blockchain {
	genesis := NewBlock(
		0,
		"DCOIN TESTNET GENESIS",
		"0",
		"GENESIS",
	)

	return &Blockchain{
		Blocks: []Block{genesis},
	}
}

func (bc *Blockchain) AddBlock(
	data string,
	validator string,
) (Block, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if data == "" || validator == "" {
		return Block{}, errors.New("invalid block data or validator")
	}

	previous := bc.Blocks[len(bc.Blocks)-1]

	block := NewBlock(
		previous.Index+1,
		data,
		previous.Hash,
		validator,
	)

	bc.Blocks = append(bc.Blocks, block)

	return block, nil
}

func (bc *Blockchain) Validate() bool {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if len(bc.Blocks) == 0 {
		return false
	}

	for i, block := range bc.Blocks {
		if !block.IsValid() {
			return false
		}

		if i == 0 {
			if block.Index != 0 ||
				block.PreviousHash != "0" {
				return false
			}
			continue
		}

		previous := bc.Blocks[i-1]

		if block.Index != previous.Index+1 {
			return false
		}

		if block.PreviousHash != previous.Hash {
			return false
		}
	}

	return true
}

func (bc *Blockchain) LatestBlock() Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.Blocks[len(bc.Blocks)-1]
}
