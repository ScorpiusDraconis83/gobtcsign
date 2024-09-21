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

// SignParam 这是待签名的交易信息，基本上最核心的信息就是这些，通过前面的逻辑能构造出这个结构，通过这个结构即可签名，签名后即可发交易
type SignParam struct {
	MsgTx     *wire.MsgTx // 既是参数也是返回值：输入时签名前的交易，而最终返回也是在这里，会得到签名后的交易
	PkScripts [][]byte
	Amounts   []int64
	NetParams *chaincfg.Params
}

// Sign 根据钱包地址和钱包私钥签名
func Sign(senderAddress string, privateKeyHex string, param *SignParam) error {
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return errors.WithMessage(err, "wrong decode private key string")
	}
	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes)

	//使用的网络不同，得到的地址也不同，因此需要确认网络
	walletAddress, err := btcutil.DecodeAddress(senderAddress, param.NetParams)
	if err != nil {
		return errors.WithMessage(err, "wrong from_address")
	}
	//开发者需要知道这是，这里有4～5种类型，各有各的签名规则
	//这里只提供有限的几种签名规则，而不是全部
	switch address := walletAddress; address.(type) {
	case *btcutil.AddressPubKeyHash: //请参考 txscript.PubKeyHashTy 的签名逻辑
		//检查钱包的地址是不是压缩的，有压缩和不压缩两种格式的地址，都是可以用的
		compress, err := CheckPKHAddressIsCompress(param.NetParams, pubKey, senderAddress)
		if err != nil {
			return errors.WithMessage(err, "wrong sign check_from_address_is_compress")
		}
		//根据是否压缩选择不同的签名逻辑
		if err := SignP2PKH(param, privKey, compress); err != nil {
			return errors.WithMessage(err, "wrong sign")
		}
	default: //其它钱包类型暂不支持
		return errors.Errorf("From地址 %s 属于 %s 类型, 类型错误", address, reflect.TypeOf(address).String()) //倒是没必要支持太多的类型
	}
	return nil
}

func CheckPKHAddressIsCompress(defaultNet *chaincfg.Params, publicKey *btcec.PublicKey, senderAddress string) (bool, error) {
	for _, isCompress := range []bool{true, false} {
		var pubKeyHash []byte
		if isCompress {
			pubKeyHash = btcutil.Hash160(publicKey.SerializeCompressed())
		} else {
			pubKeyHash = btcutil.Hash160(publicKey.SerializeUncompressed())
		}

		address, err := btcutil.NewAddressPubKeyHash(pubKeyHash, defaultNet)
		if err != nil {
			return isCompress, errors.Errorf("error=%v when is_compress=%v", err, isCompress)
		}
		if address.EncodeAddress() == senderAddress {
			return isCompress, nil
		}
	}
	return false, errors.Errorf("unknown address type. address=%s", senderAddress)
}

func SignP2PKH(signParam *SignParam, privKey *btcec.PrivateKey, compress bool) error {
	var (
		msgTx     = signParam.MsgTx
		pkScripts = signParam.PkScripts
		amounts   = signParam.Amounts
	)

	for idx := range msgTx.TxIn {
		// 使用私钥对交易输入进行签名
		// 在大多数情况下，使用压缩公钥是可以接受的，并且更常见。压缩公钥可以减小交易的大小，从而降低交易费用，并且在大多数情况下，与非压缩公钥相比，安全性没有明显的区别
		signatureScript, err := txscript.SignatureScript(msgTx, idx, pkScripts[idx], txscript.SigHashAll, privKey, compress)
		if err != nil {
			return errors.Errorf("wrong signature_script. index=%d error=%v", idx, err)
		}
		msgTx.TxIn[idx].SignatureScript = signatureScript
	}
	return VerifyP2PKHSign(msgTx, pkScripts, amounts)
}

func VerifyP2PKHSign(msgTx *wire.MsgTx, pkScripts [][]byte, amounts []int64) error {
	for idx := range msgTx.TxIn { // 这段代码的作用是创建和执行脚本引擎，用于验证指定的脚本是否有效。如果脚本验证失败，则返回错误信息。这在比特币交易的验证过程中非常重要，以确保交易的合法性和安全性。
		vm, err := txscript.NewEngine(pkScripts[idx], msgTx, idx, txscript.StandardVerifyFlags, nil, nil, amounts[idx], nil)
		if err != nil {
			return errors.Errorf("wrong vm. index=%d error=%v", idx, err)
		}
		if err = vm.Execute(); err != nil {
			return errors.Errorf("wrong vm execute. index=%d error=%v", idx, err)
		}
	}
	return nil
}

// CheckMsgTxSameWithParam 避免签名逻辑修改数量和目标位置
func CheckMsgTxSameWithParam(msgTx *wire.MsgTx, param CustomParam, netParams *chaincfg.Params) error {
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
		if seqNo := param.GetTxInSequenceNum(input); seqNo != txVin.Sequence {
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
