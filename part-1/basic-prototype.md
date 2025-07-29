# 区块链的简单模型

## 一个区块(block)
```go
type Block struct {
	Timestamp     int64 //记录创建时间
	PrevBlockHash []byte //上一个区块的Hash
	Hash          []byte //该区块的Hash(计算得出)
	Data          []byte //该区块的数据(交易等)
}
```
真实的区块链中(如BTC)，区块头(`Timestamp，PrevBlockHash, Hash`)和Data是分开实现的
```go
// BlockHeader defines information about a block and is used in the bitcoin
// block (MsgBlock) and headers (MsgHeaders) messages.
type BlockHeader struct {
    // Version of the block.  This is not the same as the protocol version.
    Version int32

    // Hash of the previous block in the block chain.
    PrevBlock chainhash.Hash

    // Merkle tree reference to hash of all transactions for the block.
    MerkleRoot chainhash.Hash

    // Time the block was created.  This is, unfortunately, encoded as a
    // uint32 on the wire and therefore is limited to 2106.
    Timestamp time.Time

    // Difficulty target for the block.
    Bits uint32

    // Nonce used to generate the block.
    Nonce uint32
}
```
> go实现的btc Header
