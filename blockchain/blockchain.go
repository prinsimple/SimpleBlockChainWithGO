// Package blockchain implements a simple blockchain and utilities.
//
// BlockChain provides methods to initialize a new blockchain, add blocks,
// iterate over the chain, and find spendable outputs.
//
// BlockChainIterator allows iterating over the blockchain.
//
// The InitBlockChain, AddBlock, FindSpendableOutputs and other functions
// implement the core blockchain logic.
package blockchain

import (
	"encoding/hex"
	// "log"
	// "errors"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First transaction from Genesis"
)

// BlockChain 区块链结构体
// LastHash 存储最后一个区块的哈希值
// DataBase 存储区块数据的数据库实例
type BlockChain struct {
	LastHash []byte
	DataBase *badger.DB
}

// BlockChainIterator 区块链迭代器
// CurrentHash 当前遍历到的区块哈希值
// DataBase 区块链数据库实例
type BlockChainIterator struct {
	CurrentHash []byte
	DataBase    *badger.DB
}

// DbExists 检查区块链数据库是否已存在
func DbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// InitBlockChain 初始化一个新的区块链
// 创建创世区块和区块链数据库
// address: 接收创世区块奖励的地址
func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	if DbExists() {
		fmt.Println("DB already exists")
		runtime.Goexit()
	}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		err = os.MkdirAll(dbPath, 0755)
		HandleError(err, "InitBlockChain")
	}

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	HandleError(err, "InitBlockChain")
	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinBaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		HandleError(err, "InitBlockChain")
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash
		return err
	})

	HandleError(err, "InitBlockChain")

	blockchain := BlockChain{lastHash, db}
	return &blockchain

}

// AddBlock 向区块链中添加新的区块
// transactions: 要打包进区块的交易列表
func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := chain.DataBase.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))

		HandleError(err, "AddBlock")
		lastHash, err = item.ValueCopy(nil)

		return err
	})
	HandleError(err, "AddBlock")

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.DataBase.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleError(err, "AddBlock")
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash
		return err
	})
	HandleError(err, "AddBlock")
}

// ContinueBlockChain 加载已存在的区块链
// 如果区块链不存在则退出程序
// address: 挖矿奖励接收地址
func ContinueBlockChain(address string) *BlockChain {
	if !DbExists() {
		fmt.Println("No existing blockchain found! Please create one")
		runtime.Goexit()
	}
	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	HandleError(err, "ContinueBlockChain")

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})
	Handle(err)
	chain := BlockChain{lastHash, db}
	return &chain
}

// Iterator 创建一个区块链迭代器
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.DataBase}

	return iter
}

// Next 获取迭代器中的下一个区块
// 返回当前哈希对应的区块，并将迭代器移动到前一个区块
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.DataBase.View(func(txn *badger.Txn) error {
		// if len(iter.CurrentHash) == 0 {
		// 	return errors.New("Error: Current Hash is empty")
		// }
		item, err := txn.Get(iter.CurrentHash)
		HandleError(err, "Next")
		encodedBlock, err := item.ValueCopy(nil)
		block = block.DeSerialize(encodedBlock)

		return err
	})
	HandleError(err, "Next")

	iter.CurrentHash = block.PrevHash

	return block
}

// FindUnspendTransactions 查找地址相关的所有未花费交易
// address: 要查询的地址
// 返回包含未花费输出的交易列表
func (chain *BlockChain) FindUnspendTransactions(address string) []Transaction {
	var unspendTxs []Transaction
	spentTXOs := make(map[string][]int)
	iter := chain.Iterator()
	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlock(address) {
					unspendTxs = append(unspendTxs, *tx)
				}
			}
			if !tx.isCoinBase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}

	}
	return unspendTxs
}

// FindUTXO 查找地址的所有未花费交易输出
// address: 要查询的地址
// 返回该地址能够使用的所有交易输出
func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspendTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlock(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

// FindSpendableOutputs 查找地址中足够支付指定金额的未花费输出
// address: 要查询的地址
// amount: 需要支付的金额
// 返回累计金额和可用的交易输出映射
func (chain *BlockChain) FindSpendableOutputs(address string, ammount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspendTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlock(address) && accumulated < ammount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= ammount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}
