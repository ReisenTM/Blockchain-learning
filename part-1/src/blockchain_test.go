package src

import (
	"fmt"
	"testing"
)

func TestBlockchain_AddBlock(t *testing.T) {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}

	// 简单断言：检查区块数
	if len(bc.blocks) != 3 { // 假设 NewBlockchain() 会加一个创世区块
		t.Errorf("expected 3 blocks, got %d", len(bc.blocks))
	}
}
