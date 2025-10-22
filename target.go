package gobtcsign

import (
	"bytes"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
)

// AddressTuple represents Bitcoin address information
// Supports both wallet address and public key script formats
// Either Address or PkScript can be provided (choose one)
// When both exist, they must match
//
// AddressTuple 代表比特币地址信息
// 支持钱包地址和公钥脚本两种格式
// Address 或 PkScript 二选一填写即可
// 当两者同时存在时，需要保证匹配
type AddressTuple struct {
	Address  string // Wallet address (choose one with PkScript) // 钱包地址（和公钥脚本二选一填写即可）
	PkScript []byte // Public key script used in tx assembly and signing // 公钥脚本（在拼装交易和签名时使用）
}

// NewAddressTuple creates AddressTuple from wallet address
// PkScript will be derived from address in subsequent logic
//
// NewAddressTuple 从钱包地址创建 AddressTuple
// PkScript 将在后续逻辑中根据地址获得
func NewAddressTuple(address string) *AddressTuple {
	return &AddressTuple{
		Address:  address,
		PkScript: nil, // Address and pk-script are mutually exclusive, pk-script derived later // address 和 pk-script 是二选一的，因此不设，在后续的逻辑里会根据地址获得 pk-script 信息
	}
}

// GetPkScript 获得公钥文本，当公钥文本存在时就用已有的，否则就根据地址计算
func (one *AddressTuple) GetPkScript(netParams *chaincfg.Params) ([]byte, error) {
	if len(one.PkScript) > 0 && len(one.Address) > 0 {
		// 这里的目的不是缓存而是两个参数都可以填，但当两个参数都填的时候就得保证匹配，避免出问题
		pkScript, err := GetAddressPkScript(one.Address, netParams)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong-address")
		}
		if !bytes.Equal(one.PkScript, pkScript) {
			return nil, errors.New("address-pk-script-mismatch")
		}
		return pkScript, nil
	}
	if len(one.PkScript) > 0 {
		return one.PkScript, nil //假如有就直接返回，否则就根据地址计算
	}
	if one.Address != "" {
		return GetAddressPkScript(one.Address, netParams) //这里不用做缓存避免增加复杂度
	}
	return nil, errors.New("no-pk-script-no-address")
}

func (one *AddressTuple) VerifyMatch(netParams *chaincfg.Params) error {
	if one.Address != "" && len(one.PkScript) > 0 {
		pkScript, err := GetAddressPkScript(one.Address, netParams)
		if err != nil {
			return errors.WithMessage(err, "wrong-address")
		}
		if !bytes.Equal(one.PkScript, pkScript) {
			return errors.New("address-pk-script-mismatch")
		}
	}
	return nil
}
