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
After the project is being built, you can run the binary with these commands:
- To add block to the chain:
```sh
./<ANY NAME> add -block "```<BLOCK DATA>```"
```
- To print the blockchain:
```sh
./<ANY NAME> print 
```
Replace ```<ANY NAME>``` with the above name you just set and ```<BLOCK DATA>``` with the data of the block 
