package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// recover 必须在 defer 函数中运行

func TestRecover(t *testing.T) {
	// recover 捕获的是祖父级调用时的异常，直接调用时无效：
	//recover()
	//panic(t.Name())

	// 直接 defer调用也无效：
	//defer recover()
	//panic(t.Name())

	// 调用时多层嵌套依然无效：
	//defer func() {
	//	func() { recover() }()
	//}()
	//panic(t.Name())

	// 必须在 defer 函数中直接调用才有效：
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, t.Name(), r)
		}
	}()
	panic(t.Name())
}
