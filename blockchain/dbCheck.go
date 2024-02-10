package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
	// "log"
)

func DbCheck() {
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	Handle(err)
	defer db.Close()
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte{}); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})
	Handle(err)
}
