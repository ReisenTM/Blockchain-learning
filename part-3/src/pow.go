package main

import (
	"blockchain-learning/part-2/src/utils"
	_ "blockchain-learning/part-2/src/utils"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

// Difficulty of mining
const targetBits = 24

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork init pow
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits)) //利用大整数表示目标值
	pow := &ProofOfWork{block: b, target: target}
	return pow
}

// nonce likes a counter
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			utils.IntToHex(pow.block.Timestamp),
			utils.IntToHex(int64(targetBits)),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run
// HashInt 是 hash 的整形表示；
// nonce 是计数器。
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte //256bit
	nonce := 0        //counter initialize
	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:]) //transform to big int
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\r%x", hash)
			//equal
			break
		} else {
			//if not equal
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	nonce := pow.block.Nonce
	data := pow.prepareData(nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
