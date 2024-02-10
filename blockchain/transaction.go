package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

type TxOutput struct {
	Value  int
	PubKey string
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]

}

func CoinBaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txInput := TxInput{[]byte{}, -1, data}
	txOutput := TxOutput{100, to}

	tx := Transaction{nil, []TxInput{txInput}, []TxOutput{txOutput}}
	tx.SetID()
	return &tx
}

func NewTransaction(from, to string, ammount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput
	acc, validOutputs := chain.FindSpendableOutputs(from, ammount)

	if acc < ammount {
		log.Panic("Error: Not enough fund")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TxOutput{ammount, to})

	if acc > ammount {
		outputs = append(outputs, TxOutput{acc - ammount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

func (tx *Transaction) isCoinBase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlock(data string) bool {
	return out.PubKey == data
}
