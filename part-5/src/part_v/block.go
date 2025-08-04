package part_v

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Transactions  []*Transaction
	Nonce         int
}

// NewBlock Create a new block
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Transactions:  transactions,
		Hash:          []byte{},
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash[:]
	return block
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	_ = gob.NewEncoder(&result).Encode(b)

	return result.Bytes()
}

// 把交易hash
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	//拼接再hash
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
