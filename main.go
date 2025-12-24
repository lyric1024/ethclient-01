package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lyric1024/ethclient-01/counter"
)

/*
1. 查询区块
编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
输出查询结果到控制台。

2. 发送交易
准备一个 Sepolia 测试网络的以太坊账户，并获取其私钥。
编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
对交易进行签名，并将签名后的交易发送到网络。
输出交易的哈希值。

3. 编写智能合约
使用 Solidity 编写一个简单的智能合约，例如一个计数器合约。
编译智能合约，生成 ABI 和字节码文件。
4. 使用 abigen 生成 Go 绑定代码
安装 abigen 工具。
使用 abigen 工具根据 ABI 和字节码文件生成 Go 绑定代码。
5. 使用生成的 Go 绑定代码与合约交互
编写 Go 代码，使用生成的 Go 绑定代码连接到 Sepolia 测试网络上的智能合约。
调用合约的方法，例如增加计数器的值。
输出调用结果。
*/

func main() {
	QueryBlock()
	TransferETH()
	CallContract()
}

// 1. 查询区块
const sepolia_url string = "https://1rpc.io/sepolia"
func QueryBlock() {
	client, err := ethclient.Dial(sepolia_url)
	if err != nil {
		log.Fatal(err)
	}

	// 实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等
	blockNumber := big.NewInt(9899209)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("区块的哈希", block.Hash().Hex())
	fmt.Println("区块的时间戳", block.Time())
	fmt.Println("区块的交易数量", block.Transactions())
}

/*
2. 构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
对交易进行签名，并将签名后的交易发送到网络。
输出交易的哈希值。
*/
func TransferETH() {
	client, err := ethclient.Dial(sepolia_url)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("生成的私钥为:", privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("get publicKeyECDSA is error")
	}
	// 获取转账地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// 转账前查询余额
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("当前地址余额为: ", balance)
	// 获取nonce值
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// suggestGasPrice
	suggestGasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	gasLimit := uint64(21000)
	// 获取chainID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 目标地址
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	value := big.NewInt(10000000000000000) // 0.01 ETH

	var data []byte
	// 生成交易
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, suggestGasPrice, data)
	// 获取私钥签名
	signerTX, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 发送交易请求
	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx hash is ", signerTX.Hash().Hex())
}


// 5. 合约的方法，例如增加计数器的值, 输出调用结果
const contractAddress = "0xBc78860D20775E4cbbA1D6Ad5e434094bF9f72ef"
const private_key = ""
func CallContract() {
	client, err := ethclient.Dial(sepolia_url)
	if err != nil {
		log.Fatal(err)
	}
	// 加载合约
	counterContract, err := counter.NewCounter(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	// 私钥
	privateKey, err := crypto.HexToECDSA(private_key)
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 创建带签名的交易发送器
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	// 调用函数
	tx, err := counterContract.Add(opt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("call Counter tx hash is ", tx.Hash().Hex())

	// 等待交易确认
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		fmt.Println("call Counter revert")
	}
	fmt.Println("call Counter blockNumber is ", receipt.BlockNumber)

	// 查询
	callOpt := &bind.CallOpts{Context: context.Background()}
	_counter, err := counterContract.GetCount(callOpt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("新值为 ", _counter)
}
