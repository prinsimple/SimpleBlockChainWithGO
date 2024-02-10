## THIS REPOSITORY WAS CREATED TO LEARN HOW TO BUILD A BLOCKCHAIN IN GO

# Go BlockChain
This is a simple blockchain written in Go

# Installation
Assuming you are using WSL.
Install Go from the official website (require 1.16 or later).
Copy the folder name go you just downloaded to : 
```sh
/usr/local/ 
```
Paste this into your bash or zsh to add enviroment path
```sh
export PATH=$PATH:/usr/local/go/bin
```
1. Clone the repository
```sh
git clone https://github.com/yourusername/golang-blockchain.git
```
2. Navigate to the project directory
```sh
cd golang-blockchain
```
3. Build the project
```sh
go build -o <ANY NAME>
```
Replace ```<ANY NAME>``` with the name you like.
# Usage
### Lesson 3
After the project is being built, you can run the binary with these commands:
- To add block to the chain:
```sh
./<ANY NAME> add -block "<BLOCK DATA>"
```
- To print the blockchain:
```sh
./<ANY NAME> print 
```
Replace ```<ANY NAME>``` with the above name you just set and ```<BLOCK DATA>``` with the data of the block 


### Lesson 4
For this lesson there are some changes to the code. So some commands are new and some are inapplicable. Here are some that works:
- To create a blockchain:
```shell
    go run main.go createblockchain -address "<ANY NAME>"
```
- To print the Blockchain:
```shell
    go run main.go print 
```
- To send some ammount of money:
```shell
go run main.go send -address "<ANY NAME>" -to "<RECEIVER>" -ammount "<AMOUNT>"
```
- To get the balance of any address:
```shell
    go run main.go getbalance -address "<ANY NAME>"
```
| Command | Description |
| --- | --- |
| `-address` | The wallet address that will be used for this operation |
| `-to`      | The recipient's wallet address |
| `-amount`   | The amount of coins to be sent |
Replace ```<ANY NAME>``` with the above name (or your address) you just set, ```<RECEIVER>``` with the address of the receiver and ```<AMOUNT>``` with the ammount of money you want to send.
