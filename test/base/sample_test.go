// See the doc: https://geektutu.com/post/hpg-benchmark.html

package base

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 运行当前 package 内的用例：go test example 或 go test .
// 运行子 package 内的用例： go test example/<package name> 或 go test ./<package name>
// 如果想递归测试当前目录下的所有的 package：go test ./... 或 go test example/...

func TestQuickStart(t *testing.T) {
	assert.Equal(t, 1, 1)
}

// go test -v -run='Pattern$'
func TestPattern(t *testing.T) {
	assert.Equal(t, 1, 1)
}

// go test -v -bench .
// b.N 表示这个用例需要运行的次数
// 对于每个用例都是不一样的
// b.N 从 1 开始，如果该用例能够在 1s 内完成，b.N 的值便会增加，再次执行
// b.N 的值大概以 1, 2, 3, 5, 10, 20, 30, 50, 100 这样的序列递增，越到后面，增加得越快

// BenchmarkFib-8           4617542               251.3 ns/op
// BenchmarkFib-8 中的 -8 即 GOMAXPROCS，默认等于 CPU 核数
// 可以通过 -cpu 参数改变 GOMAXPROCS，-cpu 支持传入一个列表作为参数
// go test -v -cpu=2,4 -bench .
func BenchmarkFib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(10)
	}
}

// go test -v -bench='Pattern$' .
func BenchmarkFibPattern(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(10)
	}
}

// go test -v -bench='Pattern$' -benchtime=3s .
// -benchtime 的值除了是时间外，还可以是具体的次数。例如，执行 30 次可以用 -benchtime=30x
// -count 参数可以用来设置 benchmark 的轮数。例如，执行 3 轮可以用 -count=3
func BenchmarkWithTime(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(10)
	}
}

func fib(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return fib(n-2) + fib(n-1)
}

// -benchmem 参数可以度量内存分配的次数
// go test -v -bench='Generate' .
// go test -v -bench='Generate' -benchmem .
func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())
	// BenchmarkGenerateWithCap-8            68          16886410 ns/op         8003644 B/op          1 allocs/op
	nums := make([]int, 0, n) // 一次分配
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func generate(n int) []int {
	rand.Seed(time.Now().UnixNano())
	// BenchmarkGenerate-8                   52          22366527 ns/op        45188576 B/op         42 allocs/op
	nums := make([]int, 0)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func BenchmarkGenerateWithCap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generateWithCap(1000000)
	}
}

func BenchmarkGenerate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generate(1000000)
	}
}

func generate2(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

// 测试不同的输入
// go test -v -bench='000$' .
func benchmarkGenerate(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		generate(i)
	}
}

// 输入变为原来的 10 倍，函数每次调用的时长也差不多是原来的 10 倍
// 说明复杂度是线性的
// BenchmarkGenerate1000
// BenchmarkGenerate1000-8            42051             28320 ns/op
// BenchmarkGenerate10000
// BenchmarkGenerate10000-8            5163            230137 ns/op
// BenchmarkGenerate100000
// BenchmarkGenerate100000-8            528           2245227 ns/op
// BenchmarkGenerate1000000
// BenchmarkGenerate1000000-8            52          22688931 ns/op

func BenchmarkGenerate1000(b *testing.B)    { benchmarkGenerate(1000, b) }
func BenchmarkGenerate10000(b *testing.B)   { benchmarkGenerate(10000, b) }
func BenchmarkGenerate100000(b *testing.B)  { benchmarkGenerate(100000, b) }
func BenchmarkGenerate1000000(b *testing.B) { benchmarkGenerate(1000000, b) }

// 如果在 benchmark 开始前，需要一些准备工作，如果准备工作比较耗时，则需要将这部分代码的耗时忽略掉
func BenchmarkFibWithTimeSleep(b *testing.B) {
	time.Sleep(time.Second * 3) // 模拟耗时准备任务
	b.ResetTimer()              // 重置定时器
	for n := 0; n < b.N; n++ {
		fib(30) // run fib(30) b.N times
	}
}

// 还有一种情况，每次函数调用前后需要一些准备工作和清理工作
// 可以使用 StopTimer 暂停计时以及使用 StartTimer 开始计时
// 例如，测试一个冒泡函数的性能，每次调用冒泡函数前，需要随机生成一个数字序列，这是非常耗时的操作
// 这种场景下，就需要使用 StopTimer 和 StartTimer 避免将这部分时间计算在内

func bubbleSort(nums []int) {
	for i := 0; i < len(nums); i++ {
		for j := 1; j < len(nums)-i; j++ {
			if nums[j] < nums[j-1] {
				nums[j], nums[j-1] = nums[j-1], nums[j]
			}
		}
	}
}

// go test -bench='Sort$' .
func BenchmarkBubbleSort(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		nums := generateWithCap(10000)
		b.StartTimer()
		bubbleSort(nums)
	}
}
