package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Difficulty 定义了挖矿难度，数值越大难度越高
const Difficulty = 12

// ProofOfWork 表示工作量证明结构
type ProofOfWork struct {
	Block  *Block   // 需要验证的区块
	Target *big.Int // 目标难度值，哈希值必须小于此目标才有效
}

// NewProof 创建一个新的工作量证明对象
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// 左移 (256-Difficulty) 位，设置目标难度
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}
	return pow
}

// InitData 准备用于哈希计算的数据
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.Block.PrevHash,           // 前一个区块的哈希
		pow.Block.HashTransaction(),  // 当前区块交易的哈希
		ToHex(int64(nonce)),          // 随机数转为十六进制
		ToHex(int64(Difficulty)),     // 难度值转为十六进制
	}, []byte{})
	return data
}

// Run 执行工作量证明算法，寻找有效的nonce值
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0 // 从0开始尝试

	// 循环直到找到有效的哈希值或达到最大整数值
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data) // 计算SHA-256哈希

		fmt.Printf("\r%x", hash) // 显示当前哈希值
		intHash.SetBytes(hash[:])

		// 如果哈希值小于目标值，说明找到了有效的nonce
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++ // 否则尝试下一个nonce值
		}
	}
	fmt.Println()

	return nonce, hash[:] // 返回有效的nonce和对应的哈希值
}

// Validate 验证区块的工作量证明是否有效
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	// 使用区块中存储的nonce重新计算哈希
	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	// 验证哈希值是否小于目标值
	return intHash.Cmp(pow.Target) == -1
}

// ToHex 将整数转换为十六进制字节数组
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
