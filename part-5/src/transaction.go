package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10 //创世区块的奖励

// 交易
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}
type TXInput struct {
	Txid      []byte //存储的是之前交易的 ID
	Vout      int    //输出在上次交易里的索引
	ScriptSig string //解锁脚本
}

// TXOutput 包含两部分
// Value: 有多少币，就是存储在 Value 里面
// ScriptPubKey: 对输出进行锁定
// 在当前实现中，ScriptPubKey 将仅用一个字符串来代替
type TXOutput struct {
	Value        int
	ScriptPubKey string //数学难题
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}

// IsCoinbase 判断是否是 coinbase 交易
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	//对交易本身取hash作ID
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// CanUnlockOutputWith
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// CanBeUnlockedWith
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// FindUnspentTransactions 找出链上与给定地址相关的所有未花费交易（UTXOs）。
// 它遍历整个区块链，记录哪些输出已被该地址花费，哪些输出仍未被花费。
// 参数：
//   - address: 要查找的地址（解锁脚本数据）
//
// 返回值：
//   - 包含该地址所有未花费输出的交易列表
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction        // 存储结果：未花费的交易
	spentTXOs := make(map[string][]int) // 记录该地址已花费的输出，key 为交易ID，value 为输出索引
	bci := bc.Iterator()                // 获取区块链迭代器，用于从最新块向前遍历

	for {
		block := bci.Next() // 获取当前区块

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID) // 当前交易ID

		Outputs: // 标签用于 continue 跳出嵌套循环
			for outIdx, out := range tx.Vout {
				// 如果该输出已经被该地址花费过，就跳过
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				// 如果这个输出属于该地址且未花费，就加入未花费交易列表
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			// 处理交易输入（除了 coinbase）
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					// 如果这个输入是该地址花费的，记录它引用的输出
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		// 如果已经遍历到创世块（PrevBlockHash为空），就退出
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// NewUTXOTransaction 普通交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		//余额不够
		log.Panic("ERROR: Not enough funds")
	}
	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)
		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) //找钱
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}

// FindSpendableOutputs 找用户可用输出(余额)
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				//追加该交易的输出索引
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}
