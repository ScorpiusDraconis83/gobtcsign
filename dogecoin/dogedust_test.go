package dogecoin

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
)

// TestNewDogeDustFee verifies DogeDustFee creation with correct configuration values
// Validates soft dust size and extra dust fee match Dogecoin network constants
//
// TestNewDogeDustFee 验证使用正确配置值创建 DogeDustFee
// 验证软灰尘大小和额外灰尘费用匹配狗狗币网络常量
func TestNewDogeDustFee(t *testing.T) {
	dustFee := NewDogeDustFee()

	require.Equal(t, btcutil.Amount(SoftDustLimit), dustFee.SoftDustSize)
	require.Equal(t, btcutil.Amount(ExtraDustsFee), dustFee.ExtraDustFee)
	require.Equal(t, btcutil.Amount(1000000), dustFee.SoftDustSize)
	require.Equal(t, btcutil.Amount(1000000), dustFee.ExtraDustFee)
}

// TestDogeDustFee_CountDustOutputs validates soft dust output counting logic
// Tests accurate identification and counting of outputs below soft dust threshold
//
// TestDogeDustFee_CountDustOutputs 验证软灰尘输出计数逻辑
// 测试准确识别和计数低于软灰尘阈值的输出
func TestDogeDustFee_CountDustOutputs(t *testing.T) {
	dustFee := NewDogeDustFee()

	outputs := []*wire.TxOut{
		wire.NewTxOut(500000, nil),  // below soft limit - dust
		wire.NewTxOut(1000000, nil), // at soft limit - not dust
		wire.NewTxOut(1500000, nil), // above soft limit - not dust
		wire.NewTxOut(800000, nil),  // below soft limit - dust
		wire.NewTxOut(999999, nil),  // below soft limit - dust
		wire.NewTxOut(2000000, nil), // above soft limit - not dust
	}

	count := dustFee.CountDustOutput(outputs)
	require.Equal(t, int64(3), count)
}

// TestDogeDustFee_SumExtraDustFee validates extra dust fee calculation accuracy
// Tests total extra fee calculation based on number of soft dust outputs
//
// TestDogeDustFee_SumExtraDustFee 验证额外灰尘费用计算准确性
// 测试基于软灰尘输出数量的总额外费用计算
func TestDogeDustFee_SumExtraDustFee(t *testing.T) {
	dustFee := NewDogeDustFee()

	outputs := []*wire.TxOut{
		wire.NewTxOut(500000, nil),  // dust
		wire.NewTxOut(1500000, nil), // not dust
		wire.NewTxOut(800000, nil),  // dust
	}

	totalFee := dustFee.SumExtraDustFee(outputs)
	require.Equal(t, btcutil.Amount(2000000), totalFee) // 2 dust * 1000000
}

// TestNewDogeDustLimit verifies DogeDustLimit checker creation and initialization
// Ensures the dust limit checker is created with proper Dogecoin network rules
//
// TestNewDogeDustLimit 验证 DogeDustLimit 检查器创建和初始化
// 确保使用正确的狗狗币网络规则创建灰尘限制检查器
func TestNewDogeDustLimit(t *testing.T) {
	dustLimit := NewDogeDustLimit()
	require.NotNil(t, dustLimit)
}

// TestDogeDustLimit_HardDust validates hard dust limit threshold enforcement
// Tests dust detection against Dogecoin hard limit across multiple output values
//
// TestDogeDustLimit_HardDust 验证硬灰尘限制阈值执行
// 测试针对狗狗币硬限制的多个输出值的灰尘检测
func TestDogeDustLimit_HardDust(t *testing.T) {
	dustLimit := NewDogeDustLimit()

	// Below hard limit - should be dust
	output := wire.NewTxOut(50000, nil)
	isDust := dustLimit.IsDustOutput(output, 0)
	require.True(t, isDust)

	// At hard limit - not dust
	output = wire.NewTxOut(100000, nil)
	isDust = dustLimit.IsDustOutput(output, 0)
	require.False(t, isDust)

	// Above hard limit - not dust
	output = wire.NewTxOut(150000, nil)
	isDust = dustLimit.IsDustOutput(output, 0)
	require.False(t, isDust)

	// Zero value - should be dust
	output = wire.NewTxOut(0, nil)
	isDust = dustLimit.IsDustOutput(output, 0)
	require.True(t, isDust)

	// Minimum valid - not dust
	output = wire.NewTxOut(100000, nil)
	isDust = dustLimit.IsDustOutput(output, 0)
	require.False(t, isDust)
}

// TestDogeDustLimit_FeeRateIndependent validates dust check independence from relay fee
// Tests that Dogecoin dust detection uses fixed threshold regardless of relay fee rate
//
// TestDogeDustLimit_FeeRateIndependent 验证灰尘检查与中继费率独立
// 测试狗狗币灰尘检测使用固定阈值而不受中继费率影响
func TestDogeDustLimit_FeeRateIndependent(t *testing.T) {
	dustLimit := NewDogeDustLimit()

	output := wire.NewTxOut(50000, nil) // below MinDustOutput

	isDust1 := dustLimit.IsDustOutput(output, 0)
	isDust2 := dustLimit.IsDustOutput(output, 1000)
	isDust3 := dustLimit.IsDustOutput(output, 10000)

	require.True(t, isDust1)
	require.True(t, isDust2)
	require.True(t, isDust3)
}

// TestDogecoin_Constants validates Dogecoin network constant values match specification
// Verifies hard dust limit, soft dust limit, and extra dust fee values are correct
//
// TestDogecoin_Constants 验证狗狗币网络常量值符合规范
// 验证硬灰尘限制、软灰尘限制和额外灰尘费用值正确
func TestDogecoin_Constants(t *testing.T) {
	require.Equal(t, int64(100000), int64(MinDustOutput))
	require.Equal(t, int64(1000000), int64(SoftDustLimit))
	require.Equal(t, int64(1000000), int64(ExtraDustsFee))
}
