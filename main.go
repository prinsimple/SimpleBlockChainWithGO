package main

import (
	// "bytes"
	// "crypto/sha256"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/tensor-programming/golang-blockchain/blockchain"
)

type CommandLine struct {
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("getbalance -address ADDRESS - get the balance of the address")
	fmt.Println("createblockchain -address ADDRESS creates a blockchain")
	// fmt.Println("add-block BLOCK_DATA -add block to the chain")
	fmt.Println("print - Prints the block in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT -Send amount from an address to another")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.DataBase.Close()
	iter := chain.Iterator()
	for {
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}

	}
}

func (cli *CommandLine) createBlockchain(address string) {
	chain := blockchain.InitBlockChain(address)
	chain.DataBase.Close()
	fmt.Println("Finished!")
}

func (cli *CommandLine) getbalance(address string) {
	chain := blockchain.ContinueBlockChain(address)
	defer chain.DataBase.Close()
	balance := 0

	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d \n", address, balance)
}

func (cli *CommandLine) send(from, to string, ammount int) {
	chain := blockchain.ContinueBlockChain(from)
	defer chain.DataBase.Close()

	tx := blockchain.NewTransaction(from, to, ammount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Send successfully!")
}

func (cli *CommandLine) run() {
	cli.validateArgs()
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The Address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.HandleError(err, "run")

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getbalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}

func main() {
	defer os.Exit(0)
	cli := CommandLine{}
	cli.run()
}
