// See the doc: https://geektutu.com/post/hpg-reduce-size.html

package main

import (
	"log"
	"net/http"
	"net/rpc"
)

type reduceCmpFSizeResult struct {
	Num, Ans int
}

type reduceCmpFSizeCalc int

// Square calculates the square of num
func (calc *reduceCmpFSizeCalc) Square(num int, result *reduceCmpFSizeResult) error {
	result.Num = num
	result.Ans = num * num
	return nil
}

func main() {
	_ = rpc.Register(new(reduceCmpFSizeCalc))
	rpc.HandleHTTP()

	log.Printf("Serving RPC server on port %d", 1234)
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("Error serving: ", err)
	}
}

// go build -o server main.go
// ls -sh server
// 8.3M server

// Go 编译器默认编译出来的程序会带有符号表和调试信息，一般来说 release 版本可以去除调试信息以减小二进制体积。
// 	-s：忽略符号表和调试信息
// 	-w：忽略 DWARFv3 调试信息，使用该选项后将无法使用gdb进行调试
// go build -ldflags="-s -w" -o server main.go
// ls -sh server
// 5.9M server

// upx 是一个常用的压缩动态库和可执行文件的工具，通常可减少 50-70% 的体积
// 	最重要的参数是压缩率，1-9，1 代表最低压缩率，9 代表最高压缩率

// 仅使用 upx
// go build -o server main.go && upx -9 server
// ls -sh server
// 4.6M server

// upx 和编译选项组合
// go build -ldflags="-s -w" -o server main.go && upx -9 server
// ls -sh server
// 2.3M server

// upx 的原理
// upx 压缩后的程序和压缩前的程序一样，无需解压仍然能够正常地运行，这种压缩方法称之为带壳压缩，压缩包含两个部分：
// 	1. 在程序开头或其他合适的地方插入解压代码；
// 	2. 将程序的其他部分压缩；
// 执行时，也包含两个部分：(也就是说，upx 在程序执行时，会有额外的解压动作，不过这个耗时几乎可以忽略)
// 	1. 首先执行的是程序开头的插入的解压代码，将原来的程序在内存中解压出来；
//	2. 再执行解压后的程序；
