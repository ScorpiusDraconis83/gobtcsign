[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yyle88/gobtcsign/release.yml?branch=main&label=BUILD)](https://github.com/yyle88/gobtcsign/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yyle88/gobtcsign)](https://pkg.go.dev/github.com/yyle88/gobtcsign)
[![Coverage Status](https://img.shields.io/coveralls/github/yyle88/gobtcsign/main.svg)](https://coveralls.io/github/yyle88/gobtcsign?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.22--1.25-lightgrey.svg)](https://github.com/yyle88/gobtcsign)
[![GitHub Release](https://img.shields.io/github/release/yyle88/gobtcsign.svg)](https://github.com/yyle88/gobtcsign/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yyle88/gobtcsign)](https://goreportcard.com/report/github.com/yyle88/gobtcsign)

---

<p align="center">
  <img
    alt="wojack-cartoon logo"
    src="assets/wojack-cartoon.jpeg"
    style="max-height: 500px; width: auto; max-width: 100%;"
  />
</p>
<h3 align="center">golang-bitcoin</h3>
<p align="center">create/sign <code>bitcoin transaction</code> with golang</p>

# gobtcsign

`gobtcsign` 简洁高效的比特币交易签名工具库，帮助开发者快速构建、签名和验证比特币交易。

`gobtcsign` 使用 golang 进行 BTC/DOGECOIN 签名，能帮助开发者探索 BTC 区块链知识。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->
## 英文文档

[ENGLISH README](README.md)
<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

---

## 安装

```bash
go get github.com/yyle88/gobtcsign
```

---

## 功能概述

以下是 `gobtcsign` 提供的核心功能：

1. **交易构建**：提供高效的交易构建工具，支持添加多个输入输出，并自动计算找零金额。通过动态手续费调整功能，用户可以灵活控制交易费用。
2. **交易大小预估**：依据输入、输出数量及脚本类型，预估交易的虚拟大小（vSize）。这有助于开发者根据实际情况设置合适的手续费率。
3. **交易签名**：兼容多种地址类型，包括 P2PKH、P2SH 和 SegWit。开发者可以使用私钥快速完成交易输入的签名。
4. **签名验证**：提供签名校验功能，确保交易签名的正确性，避免因签名问题导致交易被网络拒绝。
5. **交易序列化**：支持将签名后的交易序列化为十六进制字符串，便于直接广播至比特币网络。

---

## 依赖模块

以下是 `gobtcsign` 依赖的关键模块：

- **github.com/btcsuite/btcd**：提供比特币核心协议的实现，是构建和解析交易的基础。
- **github.com/btcsuite/btcd/btcec/v2**：用于椭圆曲线加密操作和密钥管理，支持生成和验证数字签名。
- **github.com/btcsuite/btcd/btcutil**：处理比特币地址的编码与解码操作，并提供其他常用的比特币实用工具。
- **github.com/btcsuite/btcd/chaincfg/chainhash**：提供哈希计算和链相关的常用功能。
- **github.com/btcsuite/btcwallet/wallet/txauthor**：用于构建交易的输入输出，并自动处理找零。
- **github.com/btcsuite/btcwallet/wallet/txrules**：定义比特币交易规则，包括最小手续费计算和其他限制条件。
- **github.com/btcsuite/btcwallet/wallet/txsizes**：用于计算交易的虚拟大小（vSize），便于动态调整手续费。

该项目几乎没有引用除 `github.com/btcsuite` 以外的其它包，即便如此，当您要签名交易时，依然不应该直接使用该项目，避免添加恶意代码收集您的私钥。正确的做法是fork项目，最正确的做法是拷贝代码到自己的项目里，而不要引用不可信的依赖，而且要严格审查代码，控制服务器的出入网白名单。

---

## 使用步骤

1. **初始化交易参数**：定义交易输入（UTXO）、输出目标地址及金额，同时设置 RBF（Replace-By-Fee）选项。
2. **预估交易大小与手续费**：调用库中的方法估算交易大小，依据实时费率设置合理的手续费。
3. **生成待签名交易**：根据输入的交易参数，构建待签名交易。
4. **签名交易**：使用对应私钥完成交易的数字签名。
5. **验证与序列化**：验证签名的有效性，并将交易序列化为十六进制字符串以供广播。

---

## 基本样例

### 样例1：创建比特币钱包

本样例演示如何创建 P2WPKH (SegWit) 比特币钱包，生成随机私钥并派生地址。

```go
// main 包演示 P2WPKH 钱包创建
// 生成随机私钥并派生 P2WPKH（SegWit）地址
// 输出 WIF 和十六进制格式的私钥以及比特币地址
package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func main() {
	netParams := &chaincfg.MainNetParams

	// 创建一个新的随机私钥
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		log.Fatalf("random private key error: %v", err)
	}

	// WIF（Wallet Import Format）私钥编码格式的类型
	privateWif, err := btcutil.NewWIF(privateKey, netParams, true)
	if err != nil {
		log.Fatalf("create wallet import format error: %v", err)
	}

	// 直接从私钥生成公钥
	pubKey := privateWif.PrivKey.PubKey()

	// 计算公钥哈希（P2WPKH使用的公钥哈希是公钥的SHA256和RIPEMD160哈希值）
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())

	// 创建P2WPKH地址
	witnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, netParams)
	if err != nil {
		log.Fatalf("create P2WPKH address error: %v", err)
	}

	fmt.Println("Private Key (WIF):", privateWif.String())
	fmt.Println("Private Key (Hex):", hex.EncodeToString(privateKey.Serialize()))
	fmt.Println("P2WPKH Address:", witnessPubKeyHash.EncodeAddress())
	fmt.Println("Network Name:", netParams.Name)
}
```

⬆️ **源代码：** [样例1源代码](internal/demos/demo1x/main.go)

---

### 样例2：比特币交易签名

本样例演示在测试网上使用 P2WPKH (SegWit) 地址签名比特币交易，支持 RBF。

```go
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
	// 测试网发送者地址和私钥
	// 警告：不要暴露私钥，除非准备放弃这个钱包
	const senderAddress = "tb1qvg2jksxckt96cdv9g8v9psreaggdzsrlm6arap"
	const privateKeyHex = "54bb1426611226077889d63c65f4f1fa212bcb42c2141c81e0c5409324711092"

	netParams := chaincfg.TestNet3Params

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
			{
				Target: *gobtcsign.NewAddressTuple("tb1qk0z8zhsq5hlewplv0039smnz62r2ujscz6gqjx"),
				Amount: 1234,
			},
			{
				Target: *gobtcsign.NewAddressTuple(senderAddress),
				Amount: 11855 - 11111,
			},
		},
		RBFInfo: *gobtcsign.NewRBFActive(),
	}

	// 具体费用跟实时费率以及交易体大小有关
	// 不同的交易有不同的预估值，这里省去预估过程
	mustSame(int64(11111), int64(param.GetFee()))

	// 估算交易大小（略微大于实际值）
	size, err := param.EstimateTxSize(&netParams, gobtcsign.NewNoChange())
	mustDone(err)
	fmt.Println("estimate-tx-size:", size)

	// 得到待签名的交易
	signParam, err := param.CreateTxSignParams(&netParams)
	mustDone(err)

	fmt.Println("utxo inputs:", len(signParam.InputOuts))

	// 使用私钥签名交易
	mustDone(gobtcsign.Sign(senderAddress, privateKeyHex, signParam))

	// 这是签名后的交易
	msgTx := signParam.MsgTx

	// 验证签名
	mustDone(param.VerifyMsgTxSign(msgTx, &netParams))
	// 比较信息
	mustDone(param.CheckMsgTxParam(msgTx, &netParams))

	// 获得交易哈希
	txHash := gobtcsign.GetTxHash(msgTx)
	fmt.Println("msg-tx-hash:->", txHash, "<-")
	mustSame("e587e4f65a7fa5dbba6bede6b000e8ece097671bb348db3de0e507c8b36469ad", txHash)

	// 把交易序列化得到十六进制字符串
	signedHex, err := gobtcsign.CvtMsgTxToHex(msgTx)
	mustDone(err)
	fmt.Println("raw-tx-data:->", signedHex, "<-")
	mustSame("010000000001011939727ec645869768167683487829f437cc37c664938345426d0df14e5df0e10200000000fdffffff02d204000000000000160014b3c4715e00a5ff9707ec7be2586e62d286ae4a18e80200000000000016001462152b40d8b2cbac358541d850c079ea10d1407f02483045022100e8269080acc14fd24ee13cbbdaa5ea34192f090c917b4ca3da44eda25badd58e02206813da9023bebd556a95e04e6a55c9a5fdf5dfb19746c896d7fd7f26aaa58878012102407ea64d7a9e992028a94481af95ea7d8f54870bd73e5878a014da594335ba3200000000", signedHex)

	// SendRawHexTx(txHex) - 通过这个十六进制发送交易
	// 我已经发完交易，你可以在链上看到它

	// 常见错误：
	// "-3: Amount is not a number or string" - 使用了 btcjson.NewSendRawTransactionCmd 而非 NewBitcoindSendRawTransactionCmd
	// "-26: mempool min fee not met" - 节点 minrelaytxfee 设置比较大，测试节点的费用门槛要设置小些
	fmt.Println("success")
}

// 发完交易后查发送者的账户信息：
// CONFIRMED UNSPENT: 1 OUTPUTS (0.00013089 tBTC)
// UNCONFIRMED TX COUNT: 1
// UNCONFIRMED RECEIVED: 1 OUTPUTS (0.00000744 tBTC)
// UNCONFIRMED SPENT: 1 OUTPUTS (0.00013089 tBTC)

// 发完交易后查接收者的账户信息：
// CONFIRMED UNSPENT: 1 OUTPUTS (0.00003000 tBTC)
// UNCONFIRMED TX COUNT: 1
// UNCONFIRMED RECEIVED: 1 OUTPUTS (0.00001234 tBTC)

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
```

⬆️ **源代码：** [样例2源代码](internal/demos/demo2x/main.go)

---

### 样例3：狗狗币交易签名

本样例演示在测试网上使用 P2PKH (传统) 地址签名狗狗币交易，支持 RBF。

```go
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
	// 狗狗币测试网发送者地址和私钥
	// 警告：不要暴露私钥，除非准备放弃这个钱包
	const senderAddress = "nkgVWbNrUowCG4mkWSzA7HHUDe3XyL2NaC"
	const privateKeyHex = "5f397bc72377b75db7b008a9c3fcd71651bfb138d6fc2458bb0279b9cfc8442a"

	netParams := dogecoin.TestNetParams

	// 构建狗狗币交易参数
	param := gobtcsign.BitcoinTxParams{
		VinList: []gobtcsign.VinType{
			{
				OutPoint: *gobtcsign.MustNewOutPoint(
					"173d5e1b33fc9adf64cd4b1f3b2ac73acaf0e10c967cd6fa1aa191d817d7ff77",
					3,
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
			{
				Target: *gobtcsign.NewAddressTuple(senderAddress),
				Amount: 12814705 - 222222,
			},
		},
		RBFInfo: *gobtcsign.NewRBFActive(),
	}

	// 具体费用跟实时费率以及交易体大小有关
	// 不同的交易有不同的预估值，这里省去预估过程
	mustSame(int64(222222), int64(param.GetFee()))

	// 估算交易大小（略微大于实际值）
	size, err := param.EstimateTxSize(&netParams, gobtcsign.NewNoChange())
	mustDone(err)
	fmt.Println("estimate-tx-size:", size)

	// 得到待签名的交易
	signParam, err := param.CreateTxSignParams(&netParams)
	mustDone(err)

	// 使用私钥签名交易
	mustDone(gobtcsign.Sign(senderAddress, privateKeyHex, signParam))

	// 这是签名后的交易
	msgTx := signParam.MsgTx

	// 验证签名
	mustDone(param.VerifyMsgTxSign(msgTx, &netParams))
	// 比较信息
	mustDone(param.CheckMsgTxParam(msgTx, &netParams))

	// 获得交易哈希
	txHash := gobtcsign.GetTxHash(msgTx)
	fmt.Println("msg-tx-hash:->", txHash, "<-")
	mustSame("d06f0a49c4f18e2aa520eb3bfc961602aa18c811380cb38cae3638c13883f5ed", txHash)

	// 把交易序列化得到十六进制字符串
	signedHex, err := gobtcsign.CvtMsgTxToHex(msgTx)
	mustDone(err)
	fmt.Println("raw-tx-data:->", signedHex, "<-")
	mustSame("010000000177ffd717d891a11afad67c960ce1f0ca3ac72a3b1f4bcd64df9afc331b5e3d17030000006a473044022025a41ebdb7d1a5edc5bcdb120ac339591fd95a9a084c8250a362073ffb27575202204579fa82476a52f5a28f605a827ef4866d4ba671c60363f22b523f5c27bf090a012102dfef3896f159dde1c2a972038e06ebc39c551f5f3d45e2fc9544f951fe4282f4fdffffff0287d61200000000001976a9148228d0af289894d419ddcaf6da679d8e9f0f160188ac6325c000000000001976a914b4ddb9db68061a0fec90a4bcaef21f82c8cfa1eb88ac00000000", signedHex)

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
```

⬆️ **源代码：** [样例3源代码](internal/demos/demo3x/main.go)

---

## 注意事项

1. **私钥安全性**：请勿在生产环境中暴露私钥，仅在开发或测试环境中使用演示数据。
2. **手续费设置**：根据交易大小和网络拥堵情况合理设置手续费，避免交易因手续费过低被矿工拒绝。
3. **找零地址**：在构建交易时，请确保将剩余金额转回自己的地址作为找零，以避免资金损失。
4. **网络参数**：在使用 TestNet 或 MainNet 时，请正确配置网络参数（如 `chaincfg.TestNet3Params`）。

---

通过 `gobtcsign`，开发者可以快速高效地实现比特币交易相关功能，助力区块链应用开发。

---

## 比特币入门教程

通过 `gobtcsign` 简单介绍比特币 `BTC` 的入门知识，以下是个简单的入门教程。

### 第一步-创建钱包

使用任意 **离线的代码** 创建测试钱包。 例如使用代码 [创建钱包](create_wallet_test.go)

注意不要使用在线的网页创建钱包，否则私钥容易被别人悄悄收集。

区块链的钱包创建是离线的，你能使用任意你觉得趁手的离线工具生成你的钱包（任何通过网页在线创建私钥的行为都是耍流氓）

### 第二步-找水龙头

测试币水头龙，在网上多找找总会有的，让水龙头给自己弄点测试币，这样自己就有了所谓的UTXO啦

### 第三步-尝试签名和发个交易

通过水龙头给的UTXO就可以发交易

当然实际上还是需要你具备其它能力，比如爬块技术，这样才能得到你的UTXO，否则还是不能发交易的

通过区块链浏览器 和 程序代码，你能够手动发交易，但自动化发交易还是依赖于爬块。

### 其它的-使用狗狗币学习BTC

由于狗狗币是通过LTC衍生来的，而LTC是通过BTC衍生来的，因此这个包也能用于狗狗币的签名

至于莱特币签名，没有尝试过，假如需要就试试看吧。

该包中有些狗狗币签名的样例，这是因为狗狗币的出块速度快，只几分钟就能达到6个块的确认高度，做实验或者测试相对比较便捷。
而BTC的确认达到6个块需要1小时甚至更久些，在做开发时就不太方便测试和迭代逻辑。
但BTC的资料多些，也更主流，有利于学习区块链相关的知识。
DOGE纯的模仿BTC的，逻辑99%都是互通的，因此在开发时，测试DOGE逻辑也能发现BTC的问题。
因此同时接BTC+DOGECOIN也是不错的选择。

### 特别的-注意不要遗漏找零输出

注意不要忘记找零否则将会有重大损失，详见下面的案例。

这笔交易发生在区块高度818087里面。
哈希值：b5a2af5845a8d3796308ff9840e567b14cf6bb158ff26c999e6f9a1f5448f9aa
发送方发送了139.42495946 BTC，价值5,217,651美元，而接收方仅收到了55.76998378 BTC，价值2,087,060美元。
剩余的83.65497568 BTC则是矿工费用，价值3,130,590美元。

这是一笔巨大的损失，需要特别重视，避免重蹈覆辙。

## 免责声明：

数字货币都是骗局

都是以空气币掠夺平民财富

没有公平正义可言

数字货币对中老年人是极不友好的，因为他们没有机会接触这类披着高科技外衣的割韭菜工具

数字货币对青少年也是极不友好的，因为当他们接触的时候，前面的人已经占据了大量的资源

因此妄图以数字货币，比如稍微主流的 BTC ETH TRX 代替世界货币的操作，都是不可能实现的

都不过是早先持有数字货币的八零后们的无耻幻想

扪心自问，持有几千甚至数万个比特币的人会觉得公平吗，其实不会的

因此未来还会有新事物来代替它们，而我现在也不过只是了解其中的技术，仅此而已。

该项目仅以技术学习和探索为目的而存在。

该项目作者坚定持有“坚决抵制数字货币”的立场。

---

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 许可证类型

MIT 许可证。详见 [LICENSE](LICENSE)。

---

## 🤝 项目贡献

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **发现问题？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **功能建议？** 创建 issue 讨论您的想法
- 📖 **文档疑惑？** 报告问题，帮助我们改进文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，帮助我们优化性能
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：为面向用户的更改更新文档，并使用有意义的提交消息
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来为此项目做出贡献。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/yyle88/gobtcsign.svg?variant=adaptive)](https://starchart.cc/yyle88/gobtcsign)
