
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Index        uint64 `json:"index"`
	Timestamp    int64  `json:"timestamp"`
	Data         string `json:"data"`
	PreviousHash string `json:"previous_hash"`
	Hash         string `json:"hash"`
	Validator    string `json:"validator"`
}

func NewBlock(
	index uint64,
	data string,
	previousHash string,
	validator string,
) Block {
	block := Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Data:         data,
		PreviousHash: previousHash,
		Validator:    validator,
	}

	block.Hash = block.CalculateHash()

	return block
}

func (b Block) CalculateHash() string {
	record := fmt.Sprintf(
		"%d|%d|%s|%s|%s",
		b.Index,
		b.Timestamp,
		b.Data,
		b.PreviousHash,
		b.Validator,
	)

	hash := sha256.Sum256([]byte(record))

	return hex.EncodeToString(hash[:])
}

func (b Block) IsValid() bool {
	return b.Hash == b.CalculateHash()
}
