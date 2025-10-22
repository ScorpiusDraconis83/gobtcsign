package dusts

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
)

// DustFee represents soft dust fee configuration
// Defines threshold and extra fee charged for soft dust outputs
//
// DustFee 代表软灰尘费用配置
// 定义软灰尘输出的阈值和额外收费
type DustFee struct {
	SoftDustSize btcutil.Amount // Soft dust threshold (between hard and soft limits) // 软灰尘限制（介于硬灰尘和软灰尘之间的数量）
	ExtraDustFee btcutil.Amount // Extra fee per soft dust output // 单个软灰尘额外的收费
}

// NewDustFee creates DustFee with default zero configuration
// Returns DustFee instance with soft dust disabled
//
// NewDustFee 创建使用默认零值配置的 DustFee
// 返回禁用软灰尘的 DustFee 实例
func NewDustFee() DustFee {
	return DustFee{
		SoftDustSize: 0,
		ExtraDustFee: 0,
	}
}

// CountDustOutput counts number of outputs below soft dust threshold
// Returns zero if soft dust size is not configured
//
// CountDustOutput 计算低于软灰尘阈值的输出数量
// 如果未配置软灰尘大小则返回零
func (D *DustFee) CountDustOutput(outputs []*wire.TxOut) int64 {
	var minLimit = int64(D.SoftDustSize)
	if minLimit == 0 {
		return 0
	}

	var count int64
	for _, out := range outputs {
		if out.Value < minLimit {
			count++
		}
	}
	return count
}

// SumExtraDustFee calculates total extra fee based on soft dust outputs
// Returns fee amount by multiplying dust count with extra fee rate
//
// SumExtraDustFee 根据软灰尘输出计算总额外费用
// 通过灰尘数量乘以额外费率返回费用数量
func (D *DustFee) SumExtraDustFee(outputs []*wire.TxOut) btcutil.Amount {
	return btcutil.Amount(D.CountDustOutput(outputs) * int64(D.ExtraDustFee))
}
