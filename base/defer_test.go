package base

import (
	"log"
	"os"
	"testing"
)

func TestDeferInFor(t *testing.T) {
	// defer在函数退出时才能执行，在for执行defer会导致资源延迟释放：
	// 极端情况下（比如 for 循环执行完之后程序异常），将导致所有资源没有释放
	//for i := 0; i < 5; i++ {
	//	f, err := os.Open("/path/to/file")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer f.Close()
	//}

	// 解决的方法可以在for中构造一个局部函数，在局部函数内部执行defer：
	for i := 0; i < 5; i++ {
		func() {
			f, err := os.Open("/path/to/file")
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
		}()
	}
}
