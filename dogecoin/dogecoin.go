// Package dogecoin provides Dogecoin network configuration parameters
// Registers Dogecoin MainNet, TestNet, and RegressionNet chain configs
//
// dogecoin 包提供狗狗币网络配置参数
// 注册狗狗币主网、测试网和回归测试网链配置
package dogecoin

import (
	"github.com/btcsuite/btcd/chaincfg"
)

func init() {
	if err := chaincfg.Register(&MainNetParams); err != nil {
		panic(err)
	}
	if err := chaincfg.Register(&TestNetParams); err != nil {
		panic(err)
	}
	if err := chaincfg.Register(&RegressionNetParams); err != nil {
		panic(err)
	}
}

// MainNetParams represents chain configuration for Dogecoin mainnet
// Includes network magic, address encoding, and HD wallet parameters
//
// MainNetParams 代表狗狗币主网的链配置
// 包含网络标识、地址编码和 HD 钱包参数
var MainNetParams = chaincfg.Params{
	Name: "mainnet",
	Net:  0xc0c0c0c0,

	// Address encoding magics
	PubKeyHashAddrID: 30,
	ScriptHashAddrID: 22,
	PrivateKeyID:     158,

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x02, 0xfa, 0xc3, 0x98}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x02, 0xfa, 0xca, 0xfd}, // starts with xpub

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173. Dogecoin does not actually support this, but we do not want to
	// collide with real addresses, so we specify it.
	Bech32HRPSegwit: "doge",
}

// TestNetParams represents chain configuration for Dogecoin testnet
// Includes test network magic, address encoding, and HD wallet parameters
//
// TestNetParams 代表狗狗币测试网的链配置
// 包含测试网络标识、地址编码和 HD 钱包参数
var TestNetParams = chaincfg.Params{
	Name: "testnet",
	Net:  0xfcc1b7dc,

	// Address encoding magics
	PubKeyHashAddrID: 113,
	ScriptHashAddrID: 196,
	PrivateKeyID:     241,

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with xpub

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173. Dogecoin does not actually support this, but we do not want to
	// collide with real addresses, so we specify it.
	Bech32HRPSegwit: "doget",
}

// RegressionNetParams represents chain configuration for Dogecoin regression testing network
// Uses custom network magic to avoid collision with Bitcoin RegTest
//
// RegressionNetParams 代表狗狗币回归测试网络的链配置
// 使用自定义网络标识以避免与比特币回归测试网冲突
var RegressionNetParams = chaincfg.Params{
	Name: "regtest",

	// Dogecoin has 0xdab5bffa as RegTest (same as Bitcoin's RegTest).
	// Setting it to an arbitrary value (leet_hex(dogecoin)), so that we can
	// register the regtest network.
	Net: 0xfabfb5da,

	// Address encoding magics
	PubKeyHashAddrID: 111,
	ScriptHashAddrID: 196,
	PrivateKeyID:     239,

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with xpub

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173. Dogecoin does not actually support this, but we do not want to
	// collide with real addresses, so we specify it.
	Bech32HRPSegwit: "dogert",
}
