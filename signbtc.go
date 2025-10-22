// Package gobtcsign: Bitcoin and Dogecoin transaction signing engine
// Provides comprehensive transaction building, signing, and verification capabilities
// Supports P2PKH, P2WPKH address types with auto format detection and signing
// Includes fee estimation, RBF support, and dust handling mechanisms
//
// gobtcsign: 比特币和狗狗币交易签名引擎
// 提供完整的交易构建、签名和验证功能
// 支持 P2PKH、P2WPKH 地址类型，具有自动格式检测和签名功能
// 包含费用估算、RBF 支持和灰尘处理机制
package gobtcsign

import (
	"bytes"
	"encoding/hex"
	"reflect"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

// SignParam represents the transaction information required for signing
// Contains core information needed to construct and sign transactions
// MsgTx serves as both input (unsigned) and output (signed) transaction
// InputOuts combines pkScripts and amounts from other tutorials for simplicity
//
// SignParam 代表签名所需的交易信息
// 包含构造和签名交易所需的核心信息
// MsgTx 既是输入（未签名）也是输出（已签名）的交易
// InputOuts 合并其它教程中的 pkScripts 和 amounts 以保持逻辑简洁
type SignParam struct {
	MsgTx     *wire.MsgTx      // Transaction message (input: unsigned, output: signed) // 交易消息（输入：未签名，输出：已签名）
	InputOuts []*wire.TxOut    // Combined pkScripts and amounts from inputs // 输入的 pkScripts 和 amounts 组合
	NetParams *chaincfg.Params // Network parameters (MainNet/TestNet) // 网络参数（主网/测试网）
}

// Sign signs a transaction using wallet address and private key
// Auto detects address type and applies appropriate signing method
// Supports P2WPKH (SegWit) and P2PKH (legacy) address formats
// Verifies signature after signing to ensure correctness
//
// Sign 使用钱包地址和私钥签名交易
// 自动检测地址类型并应用合适的签名方法
// 支持 P2WPKH（SegWit）和 P2PKH（传统）地址格式
// 签名后验证签名以确保正确性
func Sign(senderAddress string, privateKeyHex string, param *SignParam) error {
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return errors.WithMessage(err, "wrong decode private key string")
	}
	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes)

	// Different networks yield different addresses, so network confirmation is needed
	// 使用的网络不同，得到的地址也不同，因此需要确认网络
	walletAddress, err := btcutil.DecodeAddress(senderAddress, param.NetParams)
	if err != nil {
		return errors.WithMessage(err, "wrong from_address")
	}
	// There are 4-5 address types, each with distinct signing rules
	// This implementation provides limited but commonly used signing methods
	// 这里有4～5种地址类型，各有各的签名规则
	// 这里只提供有限的几种签名规则，而不是全部
	switch address := walletAddress; address.(type) {
	case *btcutil.AddressWitnessPubKeyHash: // txscript.WitnessV0PubKeyHashTy constant // txscript.WitnessV0PubKeyHashTy 的常量
		// Use compressed address format, uncompressed not supported
		// 这里使用压缩的地址，而不支持不压缩的
		if err := SignP2WPKH(param, privKey, true); err != nil {
			return errors.WithMessage(err, "wrong sign")
		}
	case *btcutil.AddressPubKeyHash: // Refer to txscript.PubKeyHashTy signing logic // 请参考 txscript.PubKeyHashTy 的签名逻辑
		// Check if wallet address uses compressed format (both compressed and uncompressed are valid)
		// 检查钱包的地址是不是压缩的，有压缩和不压缩两种格式的地址，都是可以用的
		compress, err := CheckPKHAddressIsCompress(param.NetParams, pubKey, senderAddress)
		if err != nil {
			return errors.WithMessage(err, "wrong sign check_from_address_is_compress")
		}
		// Choose signing logic based on compression status
		// 根据是否压缩选择不同的签名逻辑
		if err := SignP2PKH(param, privKey, compress); err != nil {
			return errors.WithMessage(err, "wrong sign")
		}
	default: // Other wallet types not yet supported (no need to support all types)
		// 其它钱包类型暂不支持（倒是没必要支持太多的类型）
		return errors.Errorf("wrong from address=%s address_type=%s not-support-this-address-type", address, reflect.TypeOf(address).String())
	}
	return nil
}

// SignP2WPKH signs SegWit (P2WPKH) transactions
// Creates witness signatures with compressed or uncompressed public keys
// Generates signature hashes and verifies signature correctness
//
// SignP2WPKH 签名 SegWit (P2WPKH) 交易
// 使用压缩或非压缩公钥创建见证签名
// 生成签名哈希并验证签名正确性
func SignP2WPKH(signParam *SignParam, privKey *btcec.PrivateKey, compress bool) error {
	var msgTx = signParam.MsgTx // Pointer pass means this serves as both parameter and return value // 这里是指针传递，因此这个既是参数也是返回值

	// Create prevOuts mapping and initialize multi-output fetcher
	// 创建 prevOuts（前置输出映射）使用 prevOuts 初始化一个多前置输出提取器
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(newPrevOutsMap(signParam))

	// Generate transaction signature hashes
	// 即可生成交易签名哈希
	sigHashes := txscript.NewTxSigHashes(msgTx, prevOutFetcher)

	// Sign each input with sigHashes
	// 接下来可以继续使用 sigHashes 进行签名
	for idx := range msgTx.TxIn {
		// Compute witness signature (P2WPKH addresses typically use compressed public keys)
		// 计算见证 P2WPKH 地址，通常使用压缩公钥
		witness, err := txscript.WitnessSignature(msgTx, sigHashes, idx, signParam.InputOuts[idx].Value, signParam.InputOuts[idx].PkScript, txscript.SigHashAll, privKey, compress)
		if err != nil {
			return errors.WithMessage(err, "witness_signature is wrong")
		}
		// Set witness data
		// 设置见证
		msgTx.TxIn[idx].Witness = witness
	}
	return VerifySign(msgTx, signParam.InputOuts, prevOutFetcher, sigHashes)
}

// VerifySign verifies transaction signature validity
// Creates and executes script engine to validate scripts
// Ensures transaction legality and security through signature verification
//
// VerifySign 验证交易签名有效性
// 创建并执行脚本引擎验证脚本
// 通过签名验证确保交易合法性和安全性
func VerifySign(msgTx *wire.MsgTx, inputOuts []*wire.TxOut, prevOutFetcher txscript.PrevOutputFetcher, sigHashes *txscript.TxSigHashes) error {
	// Setting cache size to input length is optimal (can use global cache for larger computations)
	// 设置为输入的长度是较好的，当然，更大量的计算时也可使用全局的cache
	sigCache := txscript.NewSigCache(uint(len(msgTx.TxIn)))

	// Validate inputOuts length to avoid panic and catch missing constraints after refactoring
	// 在底下的逻辑里虽然也能保证，但在这里做一次判断能避免panic，也能避免哪次重构后遗漏这个隐含的条件，因此认为在这里增加个断言还是很有必要的
	if len(inputOuts) < len(msgTx.TxIn) {
		return errors.New("wrong param-outs-length")
	}

	// Create and execute script engine to verify script validity
	// This is crucial in Bitcoin transaction verification to ensure transaction legality and security
	// 这段代码的作用是创建和执行脚本引擎，用于验证指定的脚本是否有效
	// 这在比特币交易的验证过程中非常重要，以确保交易的合法性和安全性
	for idx := range msgTx.TxIn {
		vm, err := txscript.NewEngine(inputOuts[idx].PkScript, msgTx, idx, txscript.StandardVerifyFlags, sigCache, sigHashes, inputOuts[idx].Value, prevOutFetcher)
		if err != nil {
			return errors.WithMessagef(err, "wrong new-vm-engine. index=%d", idx)
		}
		if err = vm.Execute(); err != nil {
			return errors.WithMessagef(err, "wrong check-sign-vm-execute. index=%d", idx)
		}
	}
	return nil
}

// newPrevOutsMap creates and fills previous outputs mapping
// Collects vin information and maps TxOut to corresponding OutPoint
//
// newPrevOutsMap 创建和填充前置输出映射
// 收集 vin 信息并将 TxOut 映射到对应的 OutPoint
func newPrevOutsMap(signParam *SignParam) map[wire.OutPoint]*wire.TxOut {
	var prevOutsMap = make(map[wire.OutPoint]*wire.TxOut, len(signParam.MsgTx.TxIn))
	// Collect vin information
	// 依然是只需要收集 vin 的信息
	for idx, txIn := range signParam.MsgTx.TxIn {
		// Create TxOut from amounts and pkScripts, map to corresponding OutPoint
		// 这里从 amounts 和 pkScripts 中创建 TxOut 并映射到对应的 OutPoint
		prevOutsMap[txIn.PreviousOutPoint] = wire.NewTxOut(
			signParam.InputOuts[idx].Value,
			signParam.InputOuts[idx].PkScript,
		)
	}
	return prevOutsMap
}

// CheckPKHAddressIsCompress checks if PKH address uses compressed public key format
// Tests both compressed and uncompressed formats to match sender address
// Returns compression status and error if address type is unknown
//
// CheckPKHAddressIsCompress 检查 PKH 地址是否使用压缩公钥格式
// 测试压缩和非压缩格式以匹配发送者地址
// 返回压缩状态，如果地址类型未知则返回错误
func CheckPKHAddressIsCompress(defaultNet *chaincfg.Params, publicKey *btcec.PublicKey, senderAddress string) (bool, error) {
	for _, compress := range []bool{true, false} {
		var pubKeyHash []byte
		if compress {
			pubKeyHash = btcutil.Hash160(publicKey.SerializeCompressed())
		} else {
			pubKeyHash = btcutil.Hash160(publicKey.SerializeUncompressed())
		}

		address, err := btcutil.NewAddressPubKeyHash(pubKeyHash, defaultNet)
		if err != nil {
			return compress, errors.WithMessagef(err, "wrong when address-is-compress=%v", compress)
		}
		if address.EncodeAddress() == senderAddress {
			return compress, nil
		}
	}
	return false, errors.Errorf("unknown address type. address=%s", senderAddress)
}

func SignP2PKH(signParam *SignParam, privKey *btcec.PrivateKey, compress bool) error {
	var msgTx = signParam.MsgTx // 这里是指针传递，因此这个既是参数也是返回值

	for idx := range msgTx.TxIn {
		// 使用私钥对交易输入进行签名
		// 在大多数情况下，使用压缩公钥是可以接受的，并且更常见。压缩公钥可以减小交易的大小，从而降低交易费用，并且在大多数情况下，与非压缩公钥相比，安全性没有明显的区别
		signatureScript, err := txscript.SignatureScript(msgTx, idx, signParam.InputOuts[idx].PkScript, txscript.SigHashAll, privKey, compress)
		if err != nil {
			return errors.WithMessagef(err, "wrong signature_script. index=%d", idx)
		}
		msgTx.TxIn[idx].SignatureScript = signatureScript
	}

	// 创建 prevOuts（前置输出映射） 使用 prevOuts 初始化一个多前置输出提取器
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(newPrevOutsMap(signParam))

	// 即可生成交易签名哈希
	sigHashes := txscript.NewTxSigHashes(msgTx, prevOutFetcher)

	return VerifySign(msgTx, signParam.InputOuts, prevOutFetcher, sigHashes)
}

// CheckMsgTxParam 当签完名以后最好是再用这个函数检查检查，避免签名逻辑在有BUG时修改输入或输出的内容
func (param *BitcoinTxParams) CheckMsgTxParam(msgTx *wire.MsgTx, netParams *chaincfg.Params) error {
	// 验证输入的长度是否匹配
	if len(msgTx.TxIn) != len(param.VinList) {
		return errors.Errorf("input count mismatch: got %d, expected %d", len(msgTx.TxIn), len(param.VinList))
	}
	// 验证每个输入的哈希和位置是否匹配
	for idx, txVin := range msgTx.TxIn {
		input := param.VinList[idx]
		// 检查 UTXO 的 OutPoint 是否匹配
		if txVin.PreviousOutPoint.Hash != input.OutPoint.Hash {
			return errors.Errorf("input %d outpoint-hash mismatch: got %v, expected %v", idx, txVin.PreviousOutPoint.Hash, input.OutPoint.Hash)
		}
		// 检查在交易输出中的位置是否完全匹配
		if txVin.PreviousOutPoint.Index != input.OutPoint.Index {
			return errors.Errorf("input %d outpoint-index mismatch: got %v, expected %v", idx, txVin.PreviousOutPoint.Index, input.OutPoint.Index)
		}
		// 检查 vin 的 RBF 序号是否完全匹配
		if seqNo := param.GetTxInputSequence(input); seqNo != txVin.Sequence {
			return errors.Errorf("input %d tx-in-sequence mismatch: got %v, expected %v", idx, txVin.Sequence, seqNo)
		}
	}
	// 验证输出数量是否匹配
	if len(msgTx.TxOut) != len(param.OutList) {
		return errors.Errorf("output count mismatch: got %d, expected %d", len(msgTx.TxOut), len(param.OutList))
	}
	// 验证每个输出的地址和金额是否匹配
	for idx, txVout := range msgTx.TxOut {
		output := param.OutList[idx]
		// 验证输出地址
		pkScript, err := output.Target.GetPkScript(netParams)
		if err != nil {
			return errors.Errorf("cannot get pkScript of address %s: %v", output.Target.Address, err)
		}
		if !bytes.Equal(txVout.PkScript, pkScript) {
			return errors.Errorf("output %d script mismatch: got %x, expected %x", idx, txVout.PkScript, pkScript)
		}
		// 验证输出金额
		if txVout.Value != output.Amount {
			return errors.Errorf("output %d amount mismatch: got %d, expected %d", idx, txVout.Value, output.Amount)
		}
	}
	return nil
}
