package dogecoin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMainNetParams validates Dogecoin MainNet network configuration parameters
// Verifies network IDs, address encoding values, and HD wallet key prefixes
//
// TestMainNetParams 验证狗狗币主网网络配置参数
// 验证网络 ID、地址编码值和 HD 钱包密钥前缀
func TestMainNetParams(t *testing.T) {
	require.Equal(t, "mainnet", MainNetParams.Name)
	require.EqualValues(t, 0xc0c0c0c0, MainNetParams.Net)
	require.Equal(t, uint8(30), MainNetParams.PubKeyHashAddrID)
	require.Equal(t, uint8(22), MainNetParams.ScriptHashAddrID)
	require.Equal(t, uint8(158), MainNetParams.PrivateKeyID)
	require.Equal(t, "doge", MainNetParams.Bech32HRPSegwit)

	require.Equal(t, [4]byte{0x02, 0xfa, 0xc3, 0x98}, MainNetParams.HDPrivateKeyID)
	require.Equal(t, [4]byte{0x02, 0xfa, 0xca, 0xfd}, MainNetParams.HDPublicKeyID)
}

// TestTestNetParams validates Dogecoin TestNet network configuration parameters
// Verifies test network IDs, address encoding values, and HD wallet key prefixes
//
// TestTestNetParams 验证狗狗币测试网网络配置参数
// 验证测试网络 ID、地址编码值和 HD 钱包密钥前缀
func TestTestNetParams(t *testing.T) {
	require.Equal(t, "testnet", TestNetParams.Name)
	require.EqualValues(t, 0xfcc1b7dc, TestNetParams.Net)
	require.Equal(t, uint8(113), TestNetParams.PubKeyHashAddrID)
	require.Equal(t, uint8(196), TestNetParams.ScriptHashAddrID)
	require.Equal(t, uint8(241), TestNetParams.PrivateKeyID)
	require.Equal(t, "doget", TestNetParams.Bech32HRPSegwit)

	require.Equal(t, [4]byte{0x04, 0x35, 0x83, 0x94}, TestNetParams.HDPrivateKeyID)
	require.Equal(t, [4]byte{0x04, 0x35, 0x87, 0xcf}, TestNetParams.HDPublicKeyID)
}

// TestRegressionNetParams validates Dogecoin RegressionNet configuration parameters
// Verifies regression test network IDs, address encoding, and HD wallet key prefixes
//
// TestRegressionNetParams 验证狗狗币回归测试网配置参数
// 验证回归测试网络 ID、地址编码和 HD 钱包密钥前缀
func TestRegressionNetParams(t *testing.T) {
	require.Equal(t, "regtest", RegressionNetParams.Name)
	require.EqualValues(t, 0xfabfb5da, RegressionNetParams.Net)
	require.Equal(t, uint8(111), RegressionNetParams.PubKeyHashAddrID)
	require.Equal(t, uint8(196), RegressionNetParams.ScriptHashAddrID)
	require.Equal(t, uint8(239), RegressionNetParams.PrivateKeyID)
	require.Equal(t, "dogert", RegressionNetParams.Bech32HRPSegwit)

	require.Equal(t, [4]byte{0x04, 0x35, 0x83, 0x94}, RegressionNetParams.HDPrivateKeyID)
	require.Equal(t, [4]byte{0x04, 0x35, 0x87, 0xcf}, RegressionNetParams.HDPublicKeyID)
}

// TestParams_Accessible validates network parameters exist and are accessible
// Ensures all Dogecoin network configurations initialize with valid names
//
// TestParams_Accessible 验证网络参数存在且可访问
// 确保所有狗狗币网络配置使用有效名称初始化
func TestParams_Accessible(t *testing.T) {
	require.NotNil(t, &MainNetParams)
	require.NotNil(t, &TestNetParams)
	require.NotNil(t, &RegressionNetParams)

	require.NotEmpty(t, MainNetParams.Name)
	require.NotEmpty(t, TestNetParams.Name)
	require.NotEmpty(t, RegressionNetParams.Name)
}

// TestParams_UniqueNetworkIDs validates network IDs are unique across all networks
// Ensures MainNet, TestNet, and RegressionNet use distinct network identifiers
//
// TestParams_UniqueNetworkIDs 验证所有网络的网络 ID 唯一
// 确保主网、测试网和回归测试网使用不同的网络标识符
func TestParams_UniqueNetworkIDs(t *testing.T) {
	require.NotEqual(t, MainNetParams.Net, TestNetParams.Net)
	require.NotEqual(t, MainNetParams.Net, RegressionNetParams.Net)
	require.NotEqual(t, TestNetParams.Net, RegressionNetParams.Net)
}

// TestParams_AddressEncodingDiffers validates address encoding differs across networks
// Ensures each network uses unique public key hash IDs and Bech32 prefixes
//
// TestParams_AddressEncodingDiffers 验证地址编码在各网络间不同
// 确保每个网络使用唯一的公钥哈希 ID 和 Bech32 前缀
func TestParams_AddressEncodingDiffers(t *testing.T) {
	require.NotEqual(t, MainNetParams.PubKeyHashAddrID, TestNetParams.PubKeyHashAddrID)
	require.NotEqual(t, MainNetParams.PubKeyHashAddrID, RegressionNetParams.PubKeyHashAddrID)

	require.NotEqual(t, MainNetParams.Bech32HRPSegwit, TestNetParams.Bech32HRPSegwit)
	require.NotEqual(t, MainNetParams.Bech32HRPSegwit, RegressionNetParams.Bech32HRPSegwit)
}
