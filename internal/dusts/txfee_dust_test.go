package dusts

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
)

// TestNewDustFee verifies DustFee creation with zero-value default configuration
// Ensures new DustFee instances initialize with disabled soft dust and extra fee
//
// TestNewDustFee 验证使用零值默认配置创建 DustFee
// 确保新 DustFee 实例初始化时禁用软灰尘和额外费用
func TestNewDustFee(t *testing.T) {
	dustFee := NewDustFee()
	require.Equal(t, btcutil.Amount(0), dustFee.SoftDustSize)
	require.Equal(t, btcutil.Amount(0), dustFee.ExtraDustFee)
}

// TestDustFee_CountDustOutput_ZeroLimit validates counting with zero soft dust limit
// Tests that no outputs are counted as dust when soft dust size is disabled
//
// TestDustFee_CountDustOutput_ZeroLimit 验证零软灰尘限制的计数
// 测试当软灰尘大小禁用时不计算任何输出为灰尘
func TestDustFee_CountDustOutput_ZeroLimit(t *testing.T) {
	dustFee := NewDustFee()

	outputs := []*wire.TxOut{
		wire.NewTxOut(100, nil),
		wire.NewTxOut(200, nil),
		wire.NewTxOut(50, nil),
	}

	count := dustFee.CountDustOutput(outputs)
	require.Equal(t, int64(0), count)
}

// TestDustFee_CountDustOutput_WithLimit validates dust counting with configured threshold
// Tests accurate identification of outputs below configured soft dust limit
//
// TestDustFee_CountDustOutput_WithLimit 验证使用配置阈值的灰尘计数
// 测试准确识别低于配置的软灰尘限制的输出
func TestDustFee_CountDustOutput_WithLimit(t *testing.T) {
	dustFee := DustFee{
		SoftDustSize: 1000,
		ExtraDustFee: 100,
	}

	outputs := []*wire.TxOut{
		wire.NewTxOut(500, nil),  // below limit
		wire.NewTxOut(1500, nil), // above limit
		wire.NewTxOut(800, nil),  // below limit
		wire.NewTxOut(2000, nil), // above limit
		wire.NewTxOut(999, nil),  // below limit
		wire.NewTxOut(1000, nil), // exactly at limit - not dust
	}

	count := dustFee.CountDustOutput(outputs)
	require.Equal(t, int64(3), count)
}

// TestDustFee_SumExtraDustFee validates total extra fee calculation accuracy
// Tests extra fee sum calculation based on dust output count and configured rate
//
// TestDustFee_SumExtraDustFee 验证总额外费用计算准确性
// 测试基于灰尘输出数量和配置费率的额外费用总和计算
func TestDustFee_SumExtraDustFee(t *testing.T) {
	dustFee := DustFee{
		SoftDustSize: 1000,
		ExtraDustFee: 100,
	}

	outputs := []*wire.TxOut{
		wire.NewTxOut(500, nil),  // dust
		wire.NewTxOut(1500, nil), // not dust
		wire.NewTxOut(800, nil),  // dust
		wire.NewTxOut(999, nil),  // dust
	}

	totalFee := dustFee.SumExtraDustFee(outputs)
	require.Equal(t, btcutil.Amount(300), totalFee) // 3 dust outputs * 100
}

// TestDustFee_SumExtraDustFee_EmptyOutputs validates fee calculation with empty output set
// Tests that extra dust fee is zero when no outputs are provided
//
// TestDustFee_SumExtraDustFee_EmptyOutputs 验证空输出集的费用计算
// 测试当不提供输出时额外灰尘费用为零
func TestDustFee_SumExtraDustFee_EmptyOutputs(t *testing.T) {
	dustFee := DustFee{
		SoftDustSize: 1000,
		ExtraDustFee: 100,
	}

	outputs := []*wire.TxOut{}
	totalFee := dustFee.SumExtraDustFee(outputs)
	require.Equal(t, btcutil.Amount(0), totalFee)
}

// TestDustFee_SumExtraDustFee_NoDust validates fee calculation when all outputs exceed limit
// Tests that extra dust fee is zero when no outputs qualify as dust
//
// TestDustFee_SumExtraDustFee_NoDust 验证所有输出超过限制时的费用计算
// 测试当没有输出符合灰尘条件时额外灰尘费用为零
func TestDustFee_SumExtraDustFee_NoDust(t *testing.T) {
	dustFee := DustFee{
		SoftDustSize: 1000,
		ExtraDustFee: 100,
	}

	outputs := []*wire.TxOut{
		wire.NewTxOut(1500, nil),
		wire.NewTxOut(2000, nil),
		wire.NewTxOut(1000, nil), // exactly at limit
	}

	totalFee := dustFee.SumExtraDustFee(outputs)
	require.Equal(t, btcutil.Amount(0), totalFee)
}
