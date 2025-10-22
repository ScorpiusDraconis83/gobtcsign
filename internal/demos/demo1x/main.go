// Package main demonstrates P2WPKH wallet creation
// Generates random private key and derives P2WPKH (SegWit) address
// Outputs WIF and hex format private keys along with Bitcoin address
//
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

	// Generate new random private key
	// 创建一个新的随机私钥
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		log.Fatalf("random private key error: %v", err) // 随机私钥出错: %v
	}

	// Encode private key in WIF (Wallet Import Format)
	// WIF（Wallet Import Format）私钥编码格式的类型
	privateWif, err := btcutil.NewWIF(privateKey, netParams, true)
	if err != nil {
		log.Fatalf("create wallet import format error: %v", err) // 创建钱包引用格式出错: %v
	}

	// Generate public key from private key
	// 直接从私钥生成公钥
	pubKey := privateWif.PrivKey.PubKey()

	// Calculate public key hash (P2WPKH uses SHA256 + RIPEMD160 hash of public key)
	// 计算公钥哈希（P2WPKH使用的公钥哈希是公钥的SHA256和RIPEMD160哈希值）
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())

	// Create P2WPKH address
	// 创建P2WPKH地址
	witnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, netParams)
	if err != nil {
		log.Fatalf("create P2WPKH address error: %v", err) // 创建P2WPKH地址出错: %v
	}

	fmt.Println("Private Key (WIF):", privateWif.String())                        // 私钥(WIF)
	fmt.Println("Private Key (Hex):", hex.EncodeToString(privateKey.Serialize())) // 私钥(Hex)
	fmt.Println("P2WPKH Address:", witnessPubKeyHash.EncodeAddress())             // P2WPKH地址
	fmt.Println("Network Name:", netParams.Name)                                  // 地址网络名称
}
