package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"os"
)

// Block 区块结构体
// Hash 当前区块的哈希值
// Transactions 区块包含的交易列表
// PrevHash 前一个区块的哈希值
// Nonce 工作量证明的随机数
type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

// func (b *Block) DeriveHash() {
// 	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
// 	hash := sha256.Sum256(info)
// 	b.Hash = hash[:]
// }

// HashTransaction 计算区块中所有交易的哈希值
// 将所有交易ID合并后计算SHA256哈希
// 返回32字节的哈希值
func (b *Block) HashTransaction() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// CreateBlock 创建新的区块
// txs: 要打包进区块的交易列表
// prevHash: 前一个区块的哈希值
// 返回创建好的新区块
func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, txs, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Genesis 创建创世区块
// coinbase: 创世区块中的币基交易
// 返回创世区块
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

// Serialize 将区块序列化为字节数组
// 使用 gob 编码器进行序列化
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	Handle(err)
	return res.Bytes()
}

// DeSerialize 将字节数组反序列化为区块
// data: 要反序列化的字节数组
// 返回反序列化后的区块
func (b *Block) DeSerialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	Handle(err)
	return &block
}

// Handle 处理错误
// 如果有错误则触发 panic
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// HandleError 处理错误并输出发生错误的函数名和目录
// err: 要处理的错误
// funcName: 发生错误的函数名
func HandleError(err error, funcName string) {
	if err != nil {
		dir, errDir := os.Getwd()
		if errDir != nil {
			log.Panicf("Error getting directory: %s", errDir)
		} else {
			log.Panicf("Error occurred in function %s at: %s. Error: %s", funcName, dir, err)
		}
	}
}
