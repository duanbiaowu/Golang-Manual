// See the doc: https://geektutu.com/post/hpg-escape-analysis.html

package main

import (
	"fmt"
	"math/rand"
)

type Demo struct {
	name string
}

func createDemo(name string) *Demo {
	d := new(Demo) // 局部变量 d 逃逸到堆
	d.name = name
	return d
}

func main() {
	demo := createDemo("demo")
	fmt.Println(demo)

	generate8191()
	generate8192()
	generate(1)

	in := Increase()
	fmt.Println(in()) // 1
	fmt.Println(in()) // 2
}

// go build -gcflags=-m main_pointer.go
//./main.go:11:6: can inline createDemo
//./main.go:18:20: inlining call to createDemo
//./main.go:19:13: inlining call to fmt.Println
//./main.go:11:17: leaking param: name
//./main.go:12:10: new(Demo) escapes to heap
//./main.go:18:20: new(Demo) escapes to heap
//./main.go:19:13: []interface {}{...} does not escape
//<autogenerated>:1: leaking param content: .this

// (旧版本: 具体没有确认) 空接口 interface{} 可以表示任意的类型，如果函数参数为 interface{}，编译期间很难确定其参数的具体类型，也会发生逃逸。
// 本机 go1.17.3 linux/amd64 []interface {}{...} does not escape  没有发生逃逸

// 操作系统对内核线程使用的栈空间是有大小限制 ulimit -s 默认 8M
// Go 运行时(runtime) 尝试在 goroutine 需要的时候动态地分配栈空间，goroutine 的初始栈大小为 2 KB。
//	当 goroutine 被调度时，会绑定内核线程执行，栈空间大小也不会超过操作系统的限制。
//  超过一定大小的局部变量将逃逸到堆上，不同的 Go 版本的大小限制可能不一样。

func generate8191() {
	nums := make([]int, 8191) // < 64KB
	for i := 0; i < 8191; i++ {
		nums[i] = rand.Int()
	}
}

func generate8192() {
	nums := make([]int, 8192) // = 64KB
	for i := 0; i < 8192; i++ {
		nums[i] = rand.Int()
	}
}

func generate(n int) {
	nums := make([]int, n) // 不确定大小
	for i := 0; i < n; i++ {
		nums[i] = rand.Int()
	}
}

//./main.go:48:14: make([]int, 8191) does not escape
//./main.go:55:14: make([]int, 8192) does not escape
//./main.go:62:14: make([]int, n) escapes to heap

// Increase 闭包
func Increase() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

//./main.go:79:2: moved to heap: n
//./main.go:80:9: func literal escapes to heap

// 如何利用逃逸分析提升性能 ?
// 	1. 传值会拷贝整个对象，而传指针只会拷贝指针地址，指向的对象是同一个。传指针可以减少值的拷贝，但是会导致内存分配逃逸到堆中，
//		增加垃圾回收(GC)的负担。在对象频繁创建和删除的场景下，传递指针导致的 GC 开销可能会严重影响性能。
//  2. 一般情况下，对于需要修改原对象值，或占用内存比较大的结构体，选择传指针。对于只读的占用内存较小的结构体，直接传值能够获得更好的性能。