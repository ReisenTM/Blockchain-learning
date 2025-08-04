package part_v

import (
	"bytes"
	"encoding/gob"
	"github.com/btcsuite/btcutil/base58"
	"log"
)

// TXOutput 包含两部分
// Value: 有多少币，就是存储在 Value 里面
// ScriptPubKey: 对输出进行锁定
// 在当前实现中，ScriptPubKey 将仅用一个字符串来代替
type TXOutput struct {
	Value int
	//ScriptPubKey string //数学难题
	PubKeyHash []byte
}

func (out *TXOutput) Lock(address []byte) {
	pubkeyHash := base58.Decode(string(address))   //解码地址，得到混合的hash
	pubkeyHash = pubkeyHash[1 : len(pubkeyHash)-4] //获取真正的PKH
	out.PubKeyHash = pubkeyHash
}

// IsLockedWithKey 是否是给出的pbk锁定的
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubKeyHash)
}

// NewTXOutput Create a new TX output
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}

// TXOutputs collects TXOutput
type TXOutputs struct {
	Outputs []TXOutput
}

// Serialize serializes TXOutputs
func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeOutputs deserializes TXOutputs
func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
