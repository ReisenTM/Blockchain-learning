package src

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// Block
type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Data          []byte
}

// SetHash Hash
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp, b.Hash}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:] //转为切片
}

// NewBlock Create a new block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Data:          []byte(data),
		Hash:          []byte{},
	}
	block.SetHash()
	return block
}
