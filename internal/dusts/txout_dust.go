package dusts

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
)

// DustLimit represents dust output detection checker with custom validation logic
// Provides flexible dust detection through configurable check function
//
// DustLimit 代表使用自定义验证逻辑的灰尘输出检测检查器
// 通过可配置的检查函数提供灵活的灰尘检测
type DustLimit struct {
	check func(output *wire.TxOut, relayFeePerKb btcutil.Amount) bool // Custom dust check function // 自定义灰尘检查函数
}

// NewDustLimit creates DustLimit with custom check function
// Accepts function that determines dust status based on output and fee rate
//
// NewDustLimit 使用自定义检查函数创建 DustLimit
// 接受根据输出和费率确定灰尘状态的函数
func NewDustLimit(check func(output *wire.TxOut, relayFeePerKb btcutil.Amount) bool) *DustLimit {
	return &DustLimit{check: check}
}

// IsDustOutput checks if transaction output qualifies as dust
// Uses configured check function to determine dust status
//
// IsDustOutput 检查交易输出是否符合灰尘条件
// 使用配置的检查函数确定灰尘状态
func (D *DustLimit) IsDustOutput(output *wire.TxOut, relayFeePerKb btcutil.Amount) bool {
	return D.check(output, relayFeePerKb)
}
