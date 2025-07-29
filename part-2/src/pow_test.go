package src

import (
	"fmt"
	"strconv"
	"testing"
)

func TestProofOfWork_Run(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("This is the first block")
	bc.AddBlock("This is the second block")
	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Println()
	}

	// 简单断言：检查区块数
	for _, block := range bc.blocks {
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
