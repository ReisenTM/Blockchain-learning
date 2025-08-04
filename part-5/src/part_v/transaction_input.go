package part_v

import "bytes"

type TXInput struct {
	Txid []byte //存储的是之前交易的 ID
	Vout int    //输出在上次交易里的索引
	//ScriptSig string //解锁脚本
	PubKey    []byte
	Signature []byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(pubKeyHash)
	return bytes.Equal(lockingHash, in.PubKey)
}
