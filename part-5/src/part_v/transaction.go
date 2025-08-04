package part_v

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

const subsidy = 10 //创世区块的奖励

// 交易

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() []byte {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(tx)
	if err != nil {
		log.Panic(err)
		return nil
	}
	return buf.Bytes()
}

// Hash Generate a unique hash to identify the transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = nil
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}

	txin := TXInput{[]byte{}, -1, []byte(data), nil}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.ID = tx.Hash()

	return &tx
}

// IsCoinbase 判断是否是 coinbase 交易
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// FindUnspentTransactions 找出链上与给定地址相关的所有未花费交易（UTXOs）。
// 它遍历整个区块链，记录哪些输出已被该地址花费，哪些输出仍未被花费。
// 参数：
//   - address: 要查找的地址（解锁脚本数据）
//
// 返回值：
//   - 包含该地址所有未花费输出的交易列表
func (bc *Blockchain) FindUnspentTransactions(address []byte) []Transaction {
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
				if out.IsLockedWithKey(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			// 处理交易输入（除了 coinbase）
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					// 如果这个输入是该地址花费的，记录它引用的输出
					if in.UsesKey(address) {
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

// NewUTXOTransaction 普通交易
func NewUTXOTransaction(wallet *Wallet, to string, amount int, UTXOSet *UTXOSet) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	hasPubKey := HashPubKey(wallet.PublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(hasPubKey, amount)
	if acc < amount {
		//余额不够
		log.Panic("ERROR: Not enough funds")
	}
	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)
		for _, out := range outs {
			input := TXInput{txID, out, wallet.PublicKey, nil}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, wallet.PublicKey}) //找钱
	}
	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXOSet.Blockchain.SignTransaction(&tx, wallet.PrivateKey)
	return &tx
}

// DeserializeTransaction deserializes a transaction
func DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return transaction
}

// Sign 签名
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil
		//签名
		r, s, _ := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
	}
}

// Verify 验证签名
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}

// TrimmedCopy 令pubkey和sign为空的选择性拷贝
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	for _, in := range tx.Vin {
		input := TXInput{in.Txid, in.Vout, nil, nil}
		inputs = append(inputs, input)
	}
	for _, out := range tx.Vout {
		output := TXOutput{out.Value, out.PubKeyHash}
		outputs = append(outputs, output)
	}
	return Transaction{tx.ID, inputs, outputs}
}
