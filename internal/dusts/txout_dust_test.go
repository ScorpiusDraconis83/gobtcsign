package dusts

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
)

// TestNewDustLimit verifies DustLimit creation with custom check function
// Ensures the checker initializes with provided custom dust detection logic
//
// TestNewDustLimit 验证使用自定义检查函数创建 DustLimit
// 确保检查器使用提供的自定义灰尘检测逻辑初始化
func TestNewDustLimit(t *testing.T) {
	checkFunc := func(output *wire.TxOut, relayFeePerKb btcutil.Amount) bool {
		return output.Value < 1000
	}

	dustLimit := NewDustLimit(checkFunc)
	require.NotNil(t, dustLimit)
	require.NotNil(t, dustLimit.check)
}

// TestDustLimit_IsDustOutput_SimpleThreshold validates dust detection with simple threshold
// Tests various output values against a fixed threshold to verify correct dust classification
//
// TestDustLimit_IsDustOutput_SimpleThreshold 验证使用简单阈值的灰尘检测
// 测试各种输出值与固定阈值的对比以验证正确的灰尘分类
func TestDustLimit_IsDustOutput_SimpleThreshold(t *testing.T) {
	dustLimit := NewDustLimit(func(output *wire.TxOut, relayFeePerKb btcutil.Amount) bool {
		return output.Value < 1000
	})

	// Below threshold - should be dust
	output := wire.NewTxOut(500, nil)
	isDust := dustLimit.IsDustOutput(output, 0)
	require.True(t, isDust)

	// At threshold - not dust
	output = wire.NewTxOut(1000, nil)
	isDust = dustLimit.IsDustOutput(output, 0)
	require.False(t, isDust)

	// Above threshold - not dust
	output = wire.NewTxOut(1500, nil)
	isDust = dustLimit.IsDustOutput(output, 0)
	require.False(t, isDust)

	// Zero value - should be dust
	output = wire.NewTxOut(0, nil)
	isDust = dustLimit.IsDustOutput(output, 0)
	require.True(t, isDust)
}

// TestDustLimit_IsDustOutput_WithFeeRate validates dust detection with dynamic fee rates
// Tests dust classification considering relay fee rates in dust threshold calculation
//
// TestDustLimit_IsDustOutput_WithFeeRate 验证使用动态费率的灰尘检测
// 测试在灰尘阈值计算中考虑中继费率的灰尘分类
func TestDustLimit_IsDustOutput_WithFeeRate(t *testing.T) {
	dustLimit := NewDustLimit(func(output *wire.TxOut, relayFeePerKb btcutil.Amount) bool {
		minValue := relayFeePerKb / 10 // example: 10% of relay fee
		return btcutil.Amount(output.Value) < minValue
	})

	// Dust with low fee: 50 < 100
	output := wire.NewTxOut(50, nil)
	isDust := dustLimit.IsDustOutput(output, 1000)
	require.True(t, isDust)

	// Not dust with low fee: 150 >= 100
	output = wire.NewTxOut(150, nil)
	isDust = dustLimit.IsDustOutput(output, 1000)
	require.False(t, isDust)

	// Dust with high fee: 200 < 500
	output = wire.NewTxOut(200, nil)
	isDust = dustLimit.IsDustOutput(output, 5000)
	require.True(t, isDust)

	// Not dust with high fee: 600 >= 500
	output = wire.NewTxOut(600, nil)
	isDust = dustLimit.IsDustOutput(output, 5000)
	require.False(t, isDust)
}

// TestDustLimit_IsDustOutput_ZeroRelayFee validates dust detection with zero relay fee
// Tests dust classification edge case when relay fee rate is set to zero
//
// TestDustLimit_IsDustOutput_ZeroRelayFee 验证零中继费率的灰尘检测
// 测试中继费率设置为零时的灰尘分类边缘情况
func TestDustLimit_IsDustOutput_ZeroRelayFee(t *testing.T) {
	dustLimit := NewDustLimit(func(output *wire.TxOut, relayFeePerKb btcutil.Amount) bool {
		return relayFeePerKb == 0 && output.Value < 100
	})

	output := wire.NewTxOut(50, nil)
	isDust := dustLimit.IsDustOutput(output, 0)
	require.True(t, isDust)

	isDust = dustLimit.IsDustOutput(output, 1000)
	require.False(t, isDust)
}
