package pprof

import (
	"os"
	"runtime/pprof"
	"testing"

	"github.com/pkg/profile"
)

// $ go tool pprof -http=:9999 cpu.pprof
// $ go tool pprof cpu.pprof
func TestBubbleSort(t *testing.T) {
	file := "/tmp/" + randomString(16) + "_cpu.pprof"
	println(file)
	f, _ := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()

	err := pprof.StartCPUProfile(f)
	if err != nil {
		t.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	n := 10
	for i := 0; i < 5; i++ {
		nums := generate(n)
		bubbleSort(nums)
		n *= 10
	}
}

// 2022/05/28 19:40:37 profile: memory profiling enabled (rate 1), /tmp/profile2202996550/mem.pprof
// 2022/05/28 19:40:37 profile: memory profiling disabled, /tmp/profile2202996550/mem.pprof
// $ go tool pprof -http=:9999 /tmp/profile2202996550/mem.pprof
func TestConCat(t *testing.T) {
	defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()
	concat(100)
	concat2(100) //内存占用只有 concat() 的 10% 不到
}

// testing 支持生成 CPU、memory 和 block 的 profile 文件
// -cpuprofile=$FILE
// -memprofile=$FILE, -memprofilerate=N 调整记录速率为原来的 1/N。
// -blockprofile=$FILE

// go test -v -bench="Fib$" -run='Fib$' -cpuprofile=cpu.pprof .
// go tool pprof -text cpu.pprof
// pprof 支持多种输出格式（图片、文本、Web等），直接在命令行中运行 go tool pprof 即可看到所有支持的选项
func BenchmarkFib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(10) // run fib(10) b.N times
	}
}
