# 持久化
我们虽然实现了一个简单的POW的区块链，但是并没有实现其数据存储的特性，因此我们需要数据持久化


## 数据库选型
区块链对数据库的选择并无要求，我们的学习项目中选择使用[BoltDB](https://github.com/boltdb/bolt)
它有以下优点：
1. 非常简洁
2. 用 Go 实现
3. 不需要运行一个服务器
4. 能够允许我们构造想要的数据结构

可见它是一个非常轻量化的数据库，非常适合我们的项目

## BoltDB
BoltDB中数据都以字节数组的形式存储，因此我们需要对我们的Go结构进行序列化
，这里使用`encoding/gob`,因为其实现简单

## 数据库结构

首先我们需要考虑如何在数据库进行存储
在`Bitcoin core`的实现中
分为:
- Blocks Bucket(所有块的元数据)
- ChainState Bucket(交易数据)

由于我们还没有交易功能，只用实现blocks bucket

## 实现
我们不再需要在Blockchain结构体保存全部区块
而是改为
```go
type BlockChain struct{
	tib []byte //存储区块链的尾部区块hash
	db* Blot.DB //指向数据库指针
}
```

## 区块链的读取
当区块链逐渐变大时，我们不再期望读取全部区块链到内存中，这会造成很大的性能损耗和资源占用
因此我们要实现一个区块链迭代器
```go
type BlockChainIterator struct{
    currentHash []byte
    db          *bolt.DB //用来标记附属的区块链
}
```