package gobtcsign

import (
	"github.com/btcsuite/btcwallet/wallet/txrules"
	"github.com/yyle88/gobtcsign/internal/dusts"
)

// DustFee type alias from internal dusts package
// DustFee 来自 internal dusts 包的类型别名
type DustFee = dusts.DustFee

// NewDustFee creates empty DustFee configuration for Bitcoin
// Bitcoin has no soft dust fee, returns zero config to maintain logic consistency with Dogecoin
//
// NewDustFee 创建比特币的空 DustFee 配置
// 比特币没有软灰尘收费，返回零配置以保持与狗狗币逻辑相通
func NewDustFee() DustFee {
	return dusts.NewDustFee()
}

// DustLimit type alias from internal dusts package
// DustLimit 来自 internal dusts 包的类型别名
type DustLimit = dusts.DustLimit

// NewDustLimit creates DustLimit using Bitcoin standard dust detection rules
// Uses txrules.IsDustOutput to determine dust status based on output value and relay fee
//
// NewDustLimit 创建使用比特币标准灰尘检测规则的 DustLimit
// 使用 txrules.IsDustOutput 根据输出值和中继费确定灰尘状态
func NewDustLimit() *DustLimit {
	return dusts.NewDustLimit(txrules.IsDustOutput)
}
