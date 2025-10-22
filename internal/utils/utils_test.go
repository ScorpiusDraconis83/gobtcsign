package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMustDone_Success validates MustDone does not panic with nil error
// Ensures successful execution path when no error is present
//
// TestMustDone_Success 验证 MustDone 在 nil 错误时不发生 panic
// 确保不存在错误时的成功执行路径
func TestMustDone_Success(t *testing.T) {
	require.NotPanics(t, func() {
		MustDone(nil)
	})
}

// TestMustDone_Panic validates MustDone panics when error occurs
// Ensures panic is triggered when error is present as expected
//
// TestMustDone_Panic 验证 MustDone 在发生错误时触发 panic
// 确保存在错误时按预期触发 panic
func TestMustDone_Panic(t *testing.T) {
	require.Panics(t, func() {
		MustDone(errors.New("test error"))
	})
}

// TestMustSame_Success validates MustSame does not panic with matching values
// Tests successful execution when comparing identical values across multiple types
//
// TestMustSame_Success 验证 MustSame 在值匹配时不发生 panic
// 测试比较多种类型的相同值时的成功执行
func TestMustSame_Success(t *testing.T) {
	require.NotPanics(t, func() {
		MustSame(42, 42)
		MustSame("hello", "hello")
		MustSame(true, true)
	})
}

// TestMustSame_Panic validates MustSame panics with different values
// Tests panic behavior when comparing non-matching values across types
//
// TestMustSame_Panic 验证 MustSame 在值不同时发生 panic
// 测试比较不同类型的不匹配值时的 panic 行为
func TestMustSame_Panic(t *testing.T) {
	require.Panics(t, func() {
		MustSame(42, 43)
	})

	require.Panics(t, func() {
		MustSame("hello", "world")
	})

	require.Panics(t, func() {
		MustSame(true, false)
	})
}

// TestMustSame_DifferentTypes validates MustSame with various comparable types
// Tests successful comparison across different numeric and pointer types
//
// TestMustSame_DifferentTypes 验证 MustSame 处理各种可比较类型
// 测试不同数值和指针类型的成功比较
func TestMustSame_DifferentTypes(t *testing.T) {
	require.NotPanics(t, func() {
		MustSame(int64(100), int64(100))
		MustSame(float64(3.14), float64(3.14))
		MustSame[any](nil, nil)
	})
}
