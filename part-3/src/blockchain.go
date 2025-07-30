package main

import (
	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte
	//只读事务取尾块hash
	_ = bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	newBlock := NewBlock(data, lastHash)
	_ = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash
		return err
	})
}

// NewGenesisBlock Genesis block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// NewBlockchain Create blockchain
func NewBlockchain() *Blockchain {
	var tip []byte
	db, _ := bolt.Open(dbFile, 0600, nil)
	_ = db.Update(func(tx *bolt.Tx) error {
		//查找是否存在
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			//如果不存在，就创建
			genesis := NewGenesisBlock()
			b, _ = tx.CreateBucket([]byte(blocksBucket))
			_ = b.Put(genesis.Hash, genesis.Serialize()) //block-hash -> block 结构
			_ = b.Put([]byte("l"), genesis.Hash)         //l -> 链中最后一个块的 hash
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	bc := &Blockchain{tip, db}
	return bc

}

// BlockchainIterator 迭代器
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB //用来标记附属的区块链
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}
func (bci *BlockchainIterator) Next() *Block {
	var block *Block
	_ = bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(bci.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	//迭代的核心，实现区块接力
	bci.currentHash = block.PrevBlockHash
	return block
}
