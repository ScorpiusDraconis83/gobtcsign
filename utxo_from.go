package gobtcsign

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

// GetUtxoFromInterface defines interface to retrieve UTXO sender and amount information
// Provides abstraction over different UTXO data sources (RPC client or cache)
//
// GetUtxoFromInterface 定义检索 UTXO 发送者和数量信息的接口
// 提供不同 UTXO 数据源（RPC 客户端或缓存）的抽象
type GetUtxoFromInterface interface {
	GetUtxoFrom(utxo wire.OutPoint) (*SenderAmountUtxo, error)
}

// SenderAmountUtxoClient implements GetUtxoFromInterface using RPC client
// Fetches UTXO information from Bitcoin node via RPC calls
//
// SenderAmountUtxoClient 使用 RPC 客户端实现 GetUtxoFromInterface
// 通过 RPC 调用从比特币节点获取 UTXO 信息
type SenderAmountUtxoClient struct {
	client *rpcclient.Client // Bitcoin RPC client // 比特币 RPC 客户端
}

// NewSenderAmountUtxoClient creates SenderAmountUtxoClient with RPC client
// Returns UTXO fetcher that queries Bitcoin node
//
// NewSenderAmountUtxoClient 使用 RPC 客户端创建 SenderAmountUtxoClient
// 返回查询比特币节点的 UTXO 获取器
func NewSenderAmountUtxoClient(client *rpcclient.Client) *SenderAmountUtxoClient {
	return &SenderAmountUtxoClient{client: client}
}

// GetUtxoFrom retrieves UTXO sender and amount from Bitcoin node
// Queries previous transaction to extract output details
//
// GetUtxoFrom 从比特币节点检索 UTXO 发送者和数量
// 查询前置交易以提取输出详情
func (uc *SenderAmountUtxoClient) GetUtxoFrom(utxo wire.OutPoint) (*SenderAmountUtxo, error) {
	previousUtxoTx, err := GetRawTransaction(uc.client, utxo.Hash.String())
	if err != nil {
		return nil, errors.WithMessage(err, "get-raw-transaction")
	}
	previousOutput := previousUtxoTx.Vout[utxo.Index]

	previousAmount, err := btcutil.NewAmount(previousOutput.Value)
	if err != nil {
		return nil, errors.WithMessage(err, "get-previous-amount")
	}

	utxoFrom := NewSenderAmountUtxo(
		NewAddressTuple(previousOutput.ScriptPubKey.Address),
		int64(previousAmount),
	)
	return utxoFrom, nil
}

// SenderAmountUtxo represents UTXO sender address and amount information
// Contains essential details needed for transaction signing
//
// SenderAmountUtxo 代表 UTXO 发送者地址和数量信息
// 包含交易签名所需的关键详情
type SenderAmountUtxo struct {
	sender *AddressTuple // UTXO sender address // UTXO 发送者地址
	amount int64         // UTXO amount in satoshis // UTXO 数量（单位：聪）
}

// NewSenderAmountUtxo creates SenderAmountUtxo with sender and amount
// Returns UTXO info instance with provided details
//
// NewSenderAmountUtxo 使用发送者和数量创建 SenderAmountUtxo
// 返回包含提供详情的 UTXO 信息实例
func NewSenderAmountUtxo(sender *AddressTuple, amount int64) *SenderAmountUtxo {
	return &SenderAmountUtxo{
		sender: sender,
		amount: amount,
	}
}

// SenderAmountUtxoCache implements GetUtxoFromInterface using in-memory cache
// Provides fast UTXO lookups without network calls
//
// SenderAmountUtxoCache 使用内存缓存实现 GetUtxoFromInterface
// 提供无需网络调用的快速 UTXO 查找
type SenderAmountUtxoCache struct {
	outputUtxoMap map[wire.OutPoint]*SenderAmountUtxo // UTXO cache map // UTXO 缓存映射
}

// NewSenderAmountUtxoCache creates SenderAmountUtxoCache with UTXO map
// Returns cache-based UTXO fetcher with pre-loaded data
//
// NewSenderAmountUtxoCache 使用 UTXO 映射创建 SenderAmountUtxoCache
// 返回包含预加载数据的基于缓存的 UTXO 获取器
func NewSenderAmountUtxoCache(utxoMap map[wire.OutPoint]*SenderAmountUtxo) *SenderAmountUtxoCache {
	return &SenderAmountUtxoCache{outputUtxoMap: utxoMap}
}

// GetUtxoFrom retrieves UTXO from cache
// Returns error if UTXO not found in cache
//
// GetUtxoFrom 从缓存中检索 UTXO
// 如果缓存中未找到 UTXO 则返回错误
func (uc SenderAmountUtxoCache) GetUtxoFrom(utxo wire.OutPoint) (*SenderAmountUtxo, error) {
	utxoFrom, ok := uc.outputUtxoMap[utxo]
	if !ok {
		return nil, errors.Errorf("wrong utxo[%s:%d] not-exist-in-cache", utxo.Hash.String(), utxo.Index)
	}
	return utxoFrom, nil
}
