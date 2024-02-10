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

type BlockChain struct {
	LastHash []byte
	DataBase *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	DataBase    *badger.DB
}

func DbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// InitBlockChain initializes a new blockchain DB and genesis block
// with the given coinbase transaction. It creates the DB if it doesn't exist,
// sets the genesis block as the last hash, and returns a BlockChain instance.
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

// AddBlock adds a new block to the blockchain. It retrieves the last block
// hash from the "lh" key, creates a new block with the given transactions
// and previous hash, updates the database with the new block, and sets
// the new block hash as the "lh" key.
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

// ContinueBlockChain continues the existing blockchain if one exists. It opens the
// badger DB, gets the last block hash, creates a BlockChain object and returns it.
// If no existing chain is found it prints an error and exits.
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

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.DataBase}

	return iter
}

// Next retrieves the next block in the blockchain by fetching the block for the current hash from the database,
// setting the current hash to the previous hash of that block, and returning the deserialized block.
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

// FindUnspendTransactions returns a list of transactions associated with the
// given address that have outputs that have not yet been spent. It iterates
// through the blockchain, checks each transaction's inputs and outputs, and
// builds up a mapping of which outputs have been spent. It skips over any
// outputs that have already been spent.
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

// FindUTXO returns a list of unspent transaction outputs associated with the
// provided address by searching the blockchain. It finds all unspent transactions
// for the address, and accumulates the outputs from those transactions that can
// be unlocked by the address.
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

// FindSpendableOutputs finds and returns a set of unspent transaction outputs
// that can be used to pay the given amount from the given address.
// It searches the blockchain for transactions related to the address,
// accumulates their unspent outputs until it reaches the target amount,
// and returns the accumulated amount and a mapping of transaction IDs
// to the indices of the outputs from that transaction that should be used.
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
