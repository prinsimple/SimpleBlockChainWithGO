package main

import (
	// "bytes"
	// "crypto/sha256"
	"fmt"
	"strconv"

	"github.com/tensor-programming/golang-blockchain/blockchain"
)

func main() {
	var numb int
	var data string
	chain := blockchain.InitBlockChain()
	fmt.Println("Number of Block to generate: ")
	fmt.Scanln(&numb)
	for i := 0; i < numb; i++ {
		fmt.Println("Enter data to create block: ")
		fmt.Scanln(&data)
		chain.AddBlock(data)
	}
	fmt.Println("DA CHAIN:")
	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)

		pow := blockchain.NewProof(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
