package gobtcsign

import "github.com/btcsuite/btcd/wire"

// RBFConfig represents Replace-By-Fee configuration
// Enables transaction replacement with higher fees when original fee is too low
// Prevents transactions from being stuck in node mempool
//
// RBFConfig 代表 Replace-By-Fee 配置
// 当原始手续费过低时，可使用更高手续费重发交易
// 防止交易卡在节点的内存池里
type RBFConfig struct {
	AllowRBF bool   // Enable RBF when needed (recommended to prevent stuck transactions) // 当需要RBF时需要设置（推荐启用以防止交易被卡）
	Sequence uint32 // Sequence number for RBF mechanism (allows fee replacement) // 序列号（用于RBF机制，允许增加手续费覆盖旧交易）
}

// NewRBFConfig creates RBF config with specified sequence number
// Non-MaxTxInSequenceNum values enable RBF
//
// NewRBFConfig 创建带指定序列号的 RBF 配置
// 非 MaxTxInSequenceNum 的值启用 RBF
func NewRBFConfig(sequence uint32) *RBFConfig {
	return &RBFConfig{
		AllowRBF: sequence != wire.MaxTxInSequenceNum, // Avoid zero being mistaken as RBF disabled // 避免设置为0被误认为是不使用RBF的
		Sequence: sequence,
	}
}

// NewRBFActive creates RBF config with recommended active sequence
// Uses MaxTxInSequenceNum - 2 (BTC recommended default)
// -2 is preferred over -1 for cautious and standard compliance
//
// NewRBFActive 创建带推荐激活序列的 RBF 配置
// 使用 MaxTxInSequenceNum - 2（BTC 推荐的默认值）
// 选择 -2 而不是 -1 出于谨慎性和规范性考虑
func NewRBFActive() *RBFConfig {
	return NewRBFConfig(wire.MaxTxInSequenceNum - 2) // BTC recommended default sequence // BTC推荐的默认启用RBF的就是这个数
}

// NewRBFNotUse creates RBF config with RBF disabled
// Uses MaxTxInSequenceNum to indicate no RBF
//
// NewRBFNotUse 创建禁用 RBF 的配置
// 使用 MaxTxInSequenceNum 表示不启用 RBF
func NewRBFNotUse() *RBFConfig {
	return NewRBFConfig(wire.MaxTxInSequenceNum) // Zero values mean RBF disabled // 当两个元素都为零值时表示不启用RBF机制
}

// GetSequence returns the sequence number based on RBF config
// Returns config sequence if RBF enabled, otherwise returns MaxTxInSequenceNum
//
// GetSequence 根据 RBF 配置返回序列号
// 如果启用 RBF 返回配置的序列号，否则返回 MaxTxInSequenceNum
func (cfg *RBFConfig) GetSequence() uint32 {
	if cfg.AllowRBF || cfg.Sequence > 0 { // RBF mechanism enabled with precise logic // 启用RBF机制，精确的RBF逻辑
		return cfg.Sequence
	}
	return wire.MaxTxInSequenceNum // Zero values mean RBF disabled, use max value to indicate // 当两个元素都为零值时表示不启用RBF机制，因此这里使用默认的最大值表示不启用
}
