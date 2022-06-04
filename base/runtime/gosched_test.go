// See the doc: https://chai2010.cn/advanced-go-programming-book/appendix/appendix-a-trap.html

package runtime

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGoSchedule(t *testing.T) {
	// Goroutine是协作式抢占调度，Goroutine本身不会主动放弃CPU：
	//runtime.GOMAXPROCS(1)
	//
	//go func() {
	//	for i := 0; i < 10; i++ {
	//		fmt.Println(i)
	//	}
	//}()
	//
	//for {} // 占用CPU

	// 解决的方法是在for循环加入runtime.Gosched()调度函数：
	runtime.GOMAXPROCS(1)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
	}()

	for {
		runtime.Gosched()
	}

	// 或者是通过阻塞的方式避免CPU占用：
	//runtime.GOMAXPROCS(1)
	//
	//go func() {
	//	for i := 0; i < 10; i++ {
	//		fmt.Println(i)
	//	}
	//	os.Exit(0)
	//}()
	//
	//select{}
}
