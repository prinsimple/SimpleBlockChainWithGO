package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	DataBase *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	DataBase    *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		err = os.MkdirAll(dbPath, 0755)
		Handle(err)
	}

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	Handle(err)
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			err = txn.Set(genesis.Hash, genesis.Serialize())
			if err != nil {
				return err
			}
			err = txn.Set([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.ValueCopy(lastHash)
			return err
		}
	})
	Handle(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain

}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte
	err := chain.DataBase.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(lastHash)
		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.DataBase.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.DataBase}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block
	err := iter.DataBase.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodeBlock, err := item.ValueCopy(nil)
		block = block.DeSerialize(encodeBlock)

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash
	return block
}
