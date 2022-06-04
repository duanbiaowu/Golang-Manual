package base

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGoroutineLeak(t *testing.T) {
	// Go语言是带内存自动回收的特性，因此内存一般不会泄漏。
	// 但是Goroutine确存在泄漏的情况，同时泄漏的Goroutine引用的内存同样无法被回收。

	// 程序中后台Goroutine向管道输入自然数序列，main函数中输出序列。
	// 但是当break跳出for循环的时候，后台Goroutine就处于无法被回收的状态了。
	//ch := func() <-chan int {
	//	ch := make(chan int)
	//	go func() {
	//		for i := 0; ; i++ {
	//			ch <- i
	//		}
	//	} ()
	//	return ch
	//}()
	//
	//for v := range ch {
	//	fmt.Println(v)
	//	if v == 5 {
	//		break
	//	}
	//}

	// 通过context包来避免这个问题：
	// 当main函数在break跳出循环时，通过调用cancel()来通知后台Goroutine退出，这样就避免了Goroutine的泄漏。
	ctx, cancel := context.WithCancel(context.Background())

	ch := func(ctx context.Context) <-chan int {
		ch := make(chan int)
		go func() {
			for i := 0; ; i++ {
				select {
				case <-ctx.Done():
					return
				case ch <- i:
				}
			}
		}()
		return ch
	}(ctx)

	for v := range ch {
		fmt.Println(v)
		if v == 5 {
			cancel()
			break
		}
	}

	time.Sleep(time.Second)
	assert.Equal(t, 2, runtime.NumGoroutine())
}
