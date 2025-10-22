package dogecoin

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/yyle88/gobtcsign/internal/dusts"
)

const (
	// MinDustOutput represents hard dust limit (0.001 DOGE)
	// Outputs below this value are invalid and rejected
	// Reference: https://github.com/dogecoin/dogecoin/blob/master/doc/fee-recommendation.md
	//
	// MinDustOutput 代表硬性灰尘限制（0.001 DOGE）
	// 低于此值的输出无效且会被拒绝
	// 参考：https://github.com/dogecoin/dogecoin/blob/master/doc/fee-recommendation.md
	MinDustOutput = 100000

	// SoftDustLimit represents soft dust threshold (0.01 DOGE)
	// Outputs below this value require extra 0.01 DOGE fee per output
	// Reference: https://github.com/dogecoin/dogecoin/blob/master/doc/fee-recommendation.md
	//
	// SoftDustLimit 代表弹性灰尘限制（0.01 DOGE）
	// 低于此值的输出需要每个输出额外支付 0.01 DOGE
	// 参考：https://github.com/dogecoin/dogecoin/blob/master/doc/fee-recommendation.md
	SoftDustLimit = 1000000

	// ExtraDustsFee represents extra fee charged per soft dust output (0.01 DOGE)
	// Applied at all txrules.FeeForSerializeSize calls to account for soft dust
	//
	// ExtraDustsFee 代表每个软灰尘输出额外收取的费用（0.01 DOGE）
	// 应用于所有 txrules.FeeForSerializeSize 调用以计入软灰尘费用
	ExtraDustsFee = 1000000
)

// DustFee type alias from internal dusts package
// DustFee 来自 internal dusts 包的类型别名
type DustFee = dusts.DustFee

// NewDogeDustFee creates DustFee configuration with Dogecoin soft dust rules
// Configures soft dust threshold and extra fee based on Dogecoin specification
// Reference: https://github.com/dogecoin/dogecoin/blob/b4a5d2bef20f5cca54d9c14ca118dec259e47bb4/doc/fee-recommendation.md
// Dogecoin defines soft and hard dust limits - hard dust rejected, soft dust charged extra fee
//
// NewDogeDustFee 创建使用狗狗币软灰尘规则的 DustFee 配置
// 根据狗狗币规范配置软灰尘阈值和额外费用
// 参考：https://github.com/dogecoin/dogecoin/blob/b4a5d2bef20f5cca54d9c14ca118dec259e47bb4/doc/fee-recommendation.md
// 狗狗币定义软灰尘和硬灰尘限制 - 硬灰尘会被拒绝，软灰尘会收取额外费用
func NewDogeDustFee() DustFee {
	res := dusts.NewDustFee()
	res.SoftDustSize = SoftDustLimit
	res.ExtraDustFee = ExtraDustsFee
	return res
}

// DustLimit type alias from internal dusts package
// DustLimit 来自 internal dusts 包的类型别名
type DustLimit = dusts.DustLimit

// NewDogeDustLimit creates DustLimit with Dogecoin hard dust detection rules
// Uses simple constant comparison independent of fee rate
// Outputs below MinDustOutput are considered dust and rejected
//
// NewDogeDustLimit 创建使用狗狗币硬灰尘检测规则的 DustLimit
// 使用独立于费率的简单常量比较
// 低于 MinDustOutput 的输出被视为灰尘并拒绝
func NewDogeDustLimit() *DustLimit {
	return dusts.NewDustLimit(func(output *wire.TxOut, relayFeePerKb btcutil.Amount) bool {
		// Dogecoin dust rules are simple - direct constant comparison, no fee rate dependency
		// 狗狗币的灰尘规定比较简单 - 直接和常量比较，不依赖于费率
		return output.Value < MinDustOutput
	})
}
