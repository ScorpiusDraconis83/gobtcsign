// Package main demonstrates Dogecoin transaction signing on TestNet
// Shows complete workflow: build transaction, sign, verify, and get hex output
// Uses P2PKH (legacy) address format with RBF support
//
// main 包演示狗狗币测试网交易签名
// 展示完整流程：构建交易、签名、验证和获取十六进制输出
// 使用 P2PKH（传统）地址格式，支持 RBF
package main

import (
	"fmt"

	"github.com/yyle88/gobtcsign"
	"github.com/yyle88/gobtcsign/dogecoin"
)

func main() {
	// Dogecoin TestNet sender address and private key
	// WARNING: Never expose private key unless wallet is disposable
	//
	// 狗狗币测试网发送者地址和私钥
	// 警告：不要暴露私钥，除非准备放弃这个钱包
	const senderAddress = "nkgVWbNrUowCG4mkWSzA7HHUDe3XyL2NaC"
	const privateKeyHex = "5f397bc72377b75db7b008a9c3fcd71651bfb138d6fc2458bb0279b9cfc8442a"

	netParams := dogecoin.TestNetParams

	// Build Dogecoin transaction parameters
	// 构建狗狗币交易参数
	param := gobtcsign.BitcoinTxParams{
		VinList: []gobtcsign.VinType{
			{
				OutPoint: *gobtcsign.MustNewOutPoint(
					"173d5e1b33fc9adf64cd4b1f3b2ac73acaf0e10c967cd6fa1aa191d817d7ff77",
					3, // UTXO index (hash + index form unique UTXO key) // 这里的位置是3，哈希和位置构成UTXO的主键
				),
				Sender:  *gobtcsign.NewAddressTuple(senderAddress),
				Amount:  14049272,
				RBFInfo: *gobtcsign.NewRBFNotUse(),
			},
		},
		OutList: []gobtcsign.OutType{
			{
				Target: *gobtcsign.NewAddressTuple("ng4P16anXNUrQw6VKHmoMW8NHsTkFBdNrn"),
				Amount: 1234567,
			},
			{ // Change output - return remaining balance to self, minus miner fee
				// 找零输出 - 把剩余余额返回给自己，减去矿工费用
				Target: *gobtcsign.NewAddressTuple(senderAddress),
				Amount: 12814705 - 222222, // Keep remainder minus assumed miner fee // 保留余额减去假设的矿工费
			},
		},
		RBFInfo: *gobtcsign.NewRBFActive(),
	}

	// Fee calculation depends on real-time rate and transaction size
	// Different transactions have different estimates, fee calculation skipped here
	//
	// 具体费用跟实时费率以及交易体大小有关
	// 不同的交易有不同的预估值，这里省去预估过程
	mustSame(int64(222222), int64(param.GetFee()))

	// Estimate transaction size (slightly larger than actual value)
	// 估算交易大小（略微大于实际值）
	size, err := param.EstimateTxSize(&netParams, gobtcsign.NewNoChange())
	mustDone(err)
	fmt.Println("estimate-tx-size:", size)

	// Create transaction ready to sign
	// 得到待签名的交易
	signParam, err := param.CreateTxSignParams(&netParams)
	mustDone(err)

	// Sign the transaction with private key
	// 使用私钥签名交易
	mustDone(gobtcsign.Sign(senderAddress, privateKeyHex, signParam))

	// Get signed transaction
	// 这是签名后的交易
	msgTx := signParam.MsgTx

	// Verify signature is valid
	// 验证签名
	mustDone(param.VerifyMsgTxSign(msgTx, &netParams))
	// Check transaction parameters match
	// 比较信息
	mustDone(param.CheckMsgTxParam(msgTx, &netParams))

	// Get transaction hash
	// 获得交易哈希
	txHash := gobtcsign.GetTxHash(msgTx)
	fmt.Println("msg-tx-hash:->", txHash, "<-")
	mustSame("d06f0a49c4f18e2aa520eb3bfc961602aa18c811380cb38cae3638c13883f5ed", txHash)

	// Serialize transaction to hex string
	// 把交易序列化得到十六进制字符串
	signedHex, err := gobtcsign.CvtMsgTxToHex(msgTx)
	mustDone(err)
	fmt.Println("raw-tx-data:->", signedHex, "<-")
	mustSame("010000000177ffd717d891a11afad67c960ce1f0ca3ac72a3b1f4bcd64df9afc331b5e3d17030000006a473044022025a41ebdb7d1a5edc5bcdb120ac339591fd95a9a084c8250a362073ffb27575202204579fa82476a52f5a28f605a827ef4866d4ba671c60363f22b523f5c27bf090a012102dfef3896f159dde1c2a972038e06ebc39c551f5f3d45e2fc9544f951fe4282f4fdffffff0287d61200000000001976a9148228d0af289894d419ddcaf6da679d8e9f0f160188ac6325c000000000001976a914b4ddb9db68061a0fec90a4bcaef21f82c8cfa1eb88ac00000000", signedHex)

	// SendRawHexTx(txHex) - Use this hex to broadcast Dogecoin transaction
	// Transaction already broadcasted, visible on chain
	//
	// SendRawHexTx(txHex) - 通过这个十六进制发送狗狗币交易
	// 我已经发完交易，你可以在链上看到它
	fmt.Println("success")
}

// mustDone panics if error occurs
func mustDone(err error) {
	if err != nil {
		panic(err)
	}
}

// mustSame compares two values and panics if different
func mustSame[T comparable](want, data T) {
	if want != data {
		fmt.Println("want:", want)
		fmt.Println("data:", data)
		panic("wrong")
	}
}
