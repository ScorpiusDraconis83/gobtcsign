package utils

import "fmt"

// MustDone panics if error occurs
// Ensures critical operations complete without errors
//
// MustDone 在发生错误时触发 panic
// 确保关键操作无错误完成
func MustDone(err error) {
	if err != nil {
		panic(err)
	}
}

// MustSame compares two values and panics if different
// Prints both values before panicking to aid debugging
//
// MustSame 比较两个值，如果不同则发生 panic
// 在 panic 前打印两个值以辅助调试
func MustSame[T comparable](want, data T) {
	if want != data {
		fmt.Println("want:", want)
		fmt.Println("data:", data)
		panic("wrong")
	}
}
