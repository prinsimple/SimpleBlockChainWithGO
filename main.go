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

	"github.com/notlongfen/SimpleBlockChainWithGO/blockchain"
)

// CommandLine 表示命令行接口结构
type CommandLine struct {
}

// printUsage 打印命令行使用说明
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("getbalance -address ADDRESS - get the balance of the address")
	fmt.Println("createblockchain -address ADDRESS creates a blockchain")
	// fmt.Println("add-block BLOCK_DATA -add block to the chain")
	fmt.Println("print - Prints the block in the chain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT -Send amount from an address to another")
}

// validateArgs 验证命令行参数
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit() // 安全退出程序
	}
}

// printChain 打印区块链中的所有区块
func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain("") // 加载现有区块链
	defer chain.DataBase.Close() // 确保数据库关闭
	iter := chain.Iterator() // 创建区块链迭代器
	
	// 遍历区块链中的所有区块
	for {
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		// 当到达创世区块时停止遍历
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

// createBlockchain 创建一个新的区块链
func (cli *CommandLine) createBlockchain(address string) {
	chain := blockchain.InitBlockChain(address) // 初始化区块链并将奖励发送到指定地址
	chain.DataBase.Close() // 关闭数据库连接
	fmt.Println("Finished!")
}

// getbalance 获取指定地址的余额
func (cli *CommandLine) getbalance(address string) {
	chain := blockchain.ContinueBlockChain(address) // 加载现有区块链
	defer chain.DataBase.Close() // 确保数据库关闭
	balance := 0

	// 查找地址的所有未花费交易输出
	UTXOs := chain.FindUTXO(address)

	// 计算总余额
	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d \n", address, balance)
}

// send 从一个地址向另一个地址发送代币
func (cli *CommandLine) send(from, to string, ammount int) {
	chain := blockchain.ContinueBlockChain(from) // 加载现有区块链
	defer chain.DataBase.Close() // 确保数据库关闭

	// 创建新交易
	tx := blockchain.NewTransaction(from, to, ammount, chain)
	// 将交易添加到新区块
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Send successfully!")
}

// run 运行命令行接口
func (cli *CommandLine) run() {
	cli.validateArgs() // 验证命令行参数
	
	// 定义各种命令的标志集
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	// 为各命令定义标志
	getBalanceAddress := getBalanceCmd.String("address", "", "The Address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	// 根据第一个参数选择命令
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

	// 处理获取余额命令
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getbalance(*getBalanceAddress)
	}

	// 处理创建区块链命令
	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	// 处理打印区块链命令
	if printChainCmd.Parsed() {
		cli.printChain()
	}

	// 处理发送代币命令
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}

// main 程序入口函数
func main() {
	defer os.Exit(0) // 确保程序正常退出
	cli := CommandLine{} // 创建命令行接口实例
	cli.run() // 运行命令行接口
}
