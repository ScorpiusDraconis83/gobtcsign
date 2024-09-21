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

// CustomParam 这是客户自定义的参数类型，表示要转入和转出的信息
type CustomParam struct {
	VinList []VinType //要转入进BTC节点的
	OutList []OutType //要从BTC节点转出的-这里面通常包含1个目标（转账）和1个自己（找零）
	RBFInfo RBFConfig //详见RBF机制，通常是需要启用RBF以免交易长期被卡的
}

type VinType struct {
	OutPoint wire.OutPoint //UTXO的主要信息
	Sender   AddressTuple  //发送者信息，钱包地址或者公钥文本，二选一填写即可
	Amount   int64         //发送数量，因为这里不是浮点数，因此很明显这里传的是聪的数量
	RBFInfo  RBFConfig     //还是RBF机制，前面的是控制整个交易的，这里控制单个UTXO的
}

type OutType struct {
	Target AddressTuple //接收者信息，钱包地址和公钥文本，二选一填写即可
	Amount int64        //聪的数量
}

type AddressTuple struct {
	Address  string //钱包地址，和 公钥文本 二选一填写即可
	PkScript []byte //PkScript（Public Key Script，公钥脚本），和 钱包地址 二选一填写即可
}

// GetPkScript 获得公钥文本，当公钥文本存在时就用已有的，否则就根据地址计算
func (one *AddressTuple) GetPkScript(netParams *chaincfg.Params) ([]byte, error) {
	if len(one.PkScript) > 0 {
		return one.PkScript, nil
	}
	return GetAddressPkScript(one.Address, netParams)
}

type RBFConfig struct {
	AllowRBF bool   //当需要RBF时需要设置，推荐启用RBF发交易，否则，当手续费过低时交易会卡在节点的内存池里
	Sequence uint32 //这是后来出的功能，RBF，使用更高手续费重发交易用的，当BTC交易发出到系统以后假如没人打包（手续费过低时），就可以增加手续费覆盖旧的交易
}

func (cfg *RBFConfig) GetSequence() uint32 {
	if cfg.AllowRBF || cfg.Sequence > 0 { //启用RBF机制，精确的RBF逻辑
		return cfg.Sequence
	}
	return wire.MaxTxInSequenceNum //默认值
}

// NewSignParam 根据用户的输入信息拼接交易
func NewSignParam(param CustomParam, netParams *chaincfg.Params) (*SignParam, error) {
	var msgTx = wire.NewMsgTx(wire.TxVersion)
	var pkScripts [][]byte
	var amounts []int64
	for _, input := range param.VinList {
		pkScript, err := input.Sender.GetPkScript(netParams)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong sender.address->pk-script")
		}
		pkScripts = append(pkScripts, pkScript)
		amounts = append(amounts, input.Amount)

		utxo := input.OutPoint
		txIn := wire.NewTxIn(wire.NewOutPoint(&utxo.Hash, uint32(utxo.Index)), nil, nil)
		if txIn.Sequence != wire.MaxTxInSequenceNum { //这里做个断言，因为我后面的逻辑都是基于默认值是它而写的
			return nil, errors.Errorf("wrong tx_in.sequence default value: %v", txIn.Sequence)
		}
		// 查看是否需要启用 RBF 机制
		if seqNo := calcTxInSequenceNum(input, param); seqNo != wire.MaxTxInSequenceNum {
			txIn.Sequence = seqNo
		}
		msgTx.AddTxIn(txIn)
	}
	for _, output := range param.OutList {
		pkScript, err := output.Target.GetPkScript(netParams)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong target.address->pk-script")
		}
		msgTx.AddTxOut(wire.NewTxOut(output.Amount, pkScript))
	}
	return &SignParam{
		MsgTx:     msgTx,
		PkScripts: pkScripts,
		Amounts:   amounts,
		NetParams: netParams,
	}, nil
}

func calcTxInSequenceNum(input VinType, param CustomParam) uint32 {
	// 当你确实是需要对每个交易单独设置RBF时，就可以在这里设置，单独设置到这个 vin 里面
	if seqNo := input.RBFInfo.GetSequence(); seqNo != wire.MaxTxInSequenceNum { //启用RBF机制，精确的RBF逻辑
		return seqNo
	}
	// 这里不设置也行，设置是为了启用 RBF 机制，设置到全部 vin 里面，当然前面的 RBF 会优先设置
	if seqNo := param.RBFInfo.GetSequence(); seqNo != wire.MaxTxInSequenceNum { //启用RBF机制，粗放的RBF逻辑
		// RBF (Replace-By-Fee) 是比特币网络中的一种机制。搜索官方的 “RBF” 即可得到你想要的知识
		// 简单来说 RBF 就是允许使用相同 utxo 发两次不同的交易，但只有其中的一笔能生效
		// 在启用 RBF 时发第二笔交易会报错，而允许重发时，发第二笔以后这两笔交易都会成为待打包状态，哪笔会打包和确认得看链上的打包情况
		// 通常，序列号设置为较高的值（如0xfffffffd），表示交易是可替换的
		// 因此，推荐的设置就是 txIn.Sequence = wire.MaxTxInSequenceNum - 2
		// 当然，设置为 0，1，2，3 也是可以的，只不过看着不太专业，推荐还是前面的 `0xfffffffd` 序列号
		// 理论上每个 txIn 都有独立的序列号，但是在业务中通常就是某个交易里的所有 txIn 使用相同的序列号，这样便于写CRUD逻辑
		return seqNo
	}
	// 当都没有设置的时候，就使用默认值就行
	return wire.MaxTxInSequenceNum
}

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

func VerifyP2PKHSign(tx *wire.MsgTx, pkScripts [][]byte, amounts []int64) error {
	for idx := range tx.TxIn { // 这段代码的作用是创建和执行脚本引擎，用于验证指定的脚本是否有效。如果脚本验证失败，则返回错误信息。这在比特币交易的验证过程中非常重要，以确保交易的合法性和安全性。
		vm, err := txscript.NewEngine(pkScripts[idx], tx, idx, txscript.StandardVerifyFlags, nil, nil, amounts[idx], nil)
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
		if seqNo := calcTxInSequenceNum(input, param); seqNo != txVin.Sequence {
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
