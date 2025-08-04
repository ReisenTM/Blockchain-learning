# 地址

## 钱包
钱包其实就是一个密钥对
```go

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

```

## 公私钥的生成
```go
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}
```
## 地址的生成
将一个公钥转换成一个 Base58 地址需要以下步骤：

1. 使用 RIPEMD160(SHA256(PubKey)) 哈希算法，取公钥并对其哈希两次

2. 给哈希加上地址生成算法版本的前缀

3. 对于第二步生成的结果，使用 SHA256(SHA256(payload)) 再哈希，计算校验和。校验和是结果哈希的前四个字节。

4. 将校验和附加到 version+PubKeyHash 的组合中。

5. 使用 Base58 对 version+PubKeyHash+checksum 组合进行编码。
```go
...
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}
...
```

