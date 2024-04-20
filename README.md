## THIS REPOSITORY WAS CREATED TO LEARN HOW TO BUILD A BLOCKCHAIN IN GO

# <span style = "color: rgb(0, 172, 215)">Go BlockChain</span> 
This is a simple blockchain written in <span style = "color: rgb(0, 172, 215)">Go</span>

Implementing Proof Of Work to validate and confim transactions.  
Please refer to the latest lesson on this README file.


# Installation
Assuming you are using WSL.
Install <span style = "color: rgb(0, 172, 215)">Go</span> from the official website (require 1.16 or later).
Copy the folder name <span style = "color: rgb(0, 172, 215)">Go</span> you just downloaded to : 
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
|
Replace ```<ANY NAME>``` with the above name (or your address) you just set, ```<RECEIVER>``` with the address of the receiver and ```<AMOUNT>``` with the ammount of money you want to send.
<span style = "color: yellow">❗ ATTENTION ⚠️ </span>
<!-- <p style="position: relative;">⚠️ Caution! There is important information below.</p> -->
<p style="position: relative; before: content: ''; width: 100%; height: 3px; background-color: #ff8000; margin-top: -5px;"></p>

Remember to delete the tmp files before running the commands. You can do it simply by deleting the ```tmp``` folder in the project directory or a more ```PROFESSIONAL``` way is to run the command ```rm -r tmp```.