// Package main demonstrates Bitcoin transaction signing on TestNet
// Shows complete workflow: build transaction, sign, verify, and get hex output
// Uses P2WPKH (SegWit) address format with RBF support
//
// main 包演示比特币测试网交易签名
// 展示完整流程：构建交易、签名、验证和获取十六进制输出
// 使用 P2WPKH（SegWit）地址格式，支持 RBF
package main

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/yyle88/gobtcsign"
)

func main() {
	// TestNet sender address and private key
	// WARNING: Never expose private key unless wallet is disposable
	//
	// 测试网发送者地址和私钥
	// 警告：不要暴露私钥，除非准备放弃这个钱包
	const senderAddress = "tb1qvg2jksxckt96cdv9g8v9psreaggdzsrlm6arap"
	const privateKeyHex = "54bb1426611226077889d63c65f4f1fa212bcb42c2141c81e0c5409324711092"

	netParams := chaincfg.TestNet3Params

	// Build transaction parameters with inputs and outputs
	// 构建包含输入和输出的交易参数
	param := gobtcsign.BitcoinTxParams{
		VinList: []gobtcsign.VinType{
			{
				OutPoint: *gobtcsign.MustNewOutPoint("e1f05d4ef10d6d4245839364c637cc37f429784883761668978645c67e723919", 2),
				Sender:   *gobtcsign.NewAddressTuple(senderAddress),
				Amount:   13089,
				RBFInfo:  *gobtcsign.NewRBFNotUse(),
			},
		},
		OutList: []gobtcsign.OutType{
			{ // Transfer to recipient address // 转给第一个地址
				Target: *gobtcsign.NewAddressTuple("tb1qk0z8zhsq5hlewplv0039smnz62r2ujscz6gqjx"),
				Amount: 1234,
			},
			{ // Change output - return remaining balance to self, minus miner fee
				// 找零输出 - 把剩余余额返回给自己，减去矿工费用
				Target: *gobtcsign.NewAddressTuple(senderAddress),
				Amount: 11855 - 11111, // Keep remainder minus assumed miner fee // 保留余额减去假设的矿工费
			},
		},
		RBFInfo: *gobtcsign.NewRBFActive(),
	}

	// Fee calculation depends on real-time rate and transaction size
	// Different transactions have different estimates, fee calculation skipped here
	//
	// 具体费用跟实时费率以及交易体大小有关
	// 不同的交易有不同的预估值，这里省去预估过程
	mustSame(int64(11111), int64(param.GetFee()))

	// Estimate transaction size (slightly larger than actual value)
	// 估算交易大小（略微大于实际值）
	size, err := param.EstimateTxSize(&netParams, gobtcsign.NewNoChange())
	mustDone(err)
	fmt.Println("estimate-tx-size:", size)

	// Create transaction ready to sign
	// 得到待签名的交易
	signParam, err := param.CreateTxSignParams(&netParams)
	mustDone(err)

	fmt.Println("utxo inputs:", len(signParam.InputOuts))

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
	mustSame("e587e4f65a7fa5dbba6bede6b000e8ece097671bb348db3de0e507c8b36469ad", txHash)

	// Serialize transaction to hex string
	// 把交易序列化得到十六进制字符串
	signedHex, err := gobtcsign.CvtMsgTxToHex(msgTx)
	mustDone(err)
	fmt.Println("raw-tx-data:->", signedHex, "<-")
	mustSame("010000000001011939727ec645869768167683487829f437cc37c664938345426d0df14e5df0e10200000000fdffffff02d204000000000000160014b3c4715e00a5ff9707ec7be2586e62d286ae4a18e80200000000000016001462152b40d8b2cbac358541d850c079ea10d1407f02483045022100e8269080acc14fd24ee13cbbdaa5ea34192f090c917b4ca3da44eda25badd58e02206813da9023bebd556a95e04e6a55c9a5fdf5dfb19746c896d7fd7f26aaa58878012102407ea64d7a9e992028a94481af95ea7d8f54870bd73e5878a014da594335ba3200000000", signedHex)

	// SendRawHexTx(txHex) - Use this hex to broadcast transaction
	// Transaction already broadcasted, visible on chain
	//
	// SendRawHexTx(txHex) - 通过这个十六进制发送交易
	// 我已经发完交易，你可以在链上看到它

	// Common errors:
	// "-3: Amount is not a number or string" - Using btcjson.NewSendRawTransactionCmd instead of NewBitcoindSendRawTransactionCmd
	// "-26: mempool min fee not met" - Node minrelaytxfee setting too high, test nodes should use lower threshold
	//
	// 常见错误：
	// "-3: Amount is not a number or string" - 使用了 btcjson.NewSendRawTransactionCmd 而非 NewBitcoindSendRawTransactionCmd
	// "-26: mempool min fee not met" - 节点 minrelaytxfee 设置比较大，测试节点的费用门槛要设置小些
	fmt.Println("success")
}

// After broadcasting transaction - sender account status:
// CONFIRMED UNSPENT: 1 OUTPUTS (0.00013089 tBTC)
// UNCONFIRMED TX COUNT: 1
// UNCONFIRMED RECEIVED: 1 OUTPUTS (0.00000744 tBTC)
// UNCONFIRMED SPENT: 1 OUTPUTS (0.00013089 tBTC)
//
// 发完交易后查发送者的账户信息：
// CONFIRMED UNSPENT: 1 OUTPUTS (0.00013089 tBTC)
// UNCONFIRMED TX COUNT: 1
// UNCONFIRMED RECEIVED: 1 OUTPUTS (0.00000744 tBTC)
// UNCONFIRMED SPENT: 1 OUTPUTS (0.00013089 tBTC)

// After broadcasting transaction - recipient account status:
// CONFIRMED UNSPENT: 1 OUTPUTS (0.00003000 tBTC)
// UNCONFIRMED TX COUNT: 1
// UNCONFIRMED RECEIVED: 1 OUTPUTS (0.00001234 tBTC)
//
// 发完交易后查接收者的账户信息：
// CONFIRMED UNSPENT: 1 OUTPUTS (0.00003000 tBTC)
// UNCONFIRMED TX COUNT: 1
// UNCONFIRMED RECEIVED: 1 OUTPUTS (0.00001234 tBTC)

// Wait for blockchain confirmation - higher fee means faster confirmation
// Otherwise wait patiently or increase fee by reconstructing transaction
//
// 接下来等待链的确认即可，给的手续费越高确认越快
// 否则就需要耐心等待，或者提高手续费重新构造和发送交易

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
