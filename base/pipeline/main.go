package main

import (
	"fmt"
	"math"
	"strconv"
	"sync"
)

// 1. Channel 管理 --------------------------------------------------
// Channel 转发函数
func echo(numbs []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range numbs {
			out <- n
		}
		close(out)
	}()
	return out
}

// 平方函数
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// 奇数过滤函数
func odd(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if n&1 == 1 {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

// 求和函数
func sum(in <-chan int, name string) <-chan int {
	out := make(chan int)
	go func() {
		res := 0
		for n := range in {
			res += n
		}
		out <- res
		fmt.Println("sum closing", name)
		close(out)
	}()
	return out
}

// 执行了类似 Unix/Linux 命令 echo $nums | odd | sq | sum
// 	nums := []int{1, 2, 3, 4, 5}
//	for n := range sum(sq(odd(echo(nums)))) {
//		fmt.Println(n)
//	}

// 如果不想多层嵌套，可以使用一个代理来完成
type EchoFunc func([]int) <-chan int
type PipeFunc func(<-chan int) <-chan int

func pipeline(nums []int, echo EchoFunc, pipeFns ...PipeFunc) <-chan int {
	ch := echo(nums)
	for i := range pipeFns {
		ch = pipeFns[i](ch)
	}
	return ch
}

//nums := []int{1, 2, 3, 4, 5}
//for n := range pipeline(nums, echo, odd, sq, sum) {
//	fmt.Println(n)
//}

// 2. Fan In/Out 一对多 OR 多对一 ----------------------------------------
// 通过并发的方式来对数组中的质数求和
// 先把数组分段求和，然后汇总
func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func isPrime(value int) bool {
	for i := 2; i <= int(math.Floor(float64(value)/2)); i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func prime(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if isPrime(n) {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

func merge(cs []<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan int) {
			for n := range c {
				out <- n
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	nums := makeRange(1, 10000)
	in := echo(nums)

	const nProcess = 5
	var chans [nProcess]<-chan int
	for i := range chans {
		chans[i] = sum(prime(in), "sum-"+strconv.Itoa(i+1))
	}

	for n := range sum(merge(chans[:]), "last sum") {
		fmt.Println(n)
	}
}
