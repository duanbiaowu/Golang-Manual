package performance

import (
	"fmt"
	"strings"
	"testing"
)

func TestBuilderConcat(t *testing.T) {
	var str = randomString(10)
	var builder strings.Builder
	cap := 0
	for i := 0; i < 10000; i++ {
		if builder.Cap() != cap {
			cap = builder.Cap()
			fmt.Printf("cap = %d\n", cap)
		}
		builder.WriteString(str)
	}
}

func benchmarkStrConcat(b *testing.B, f func(int, string) string) {
	var str = randomString(10)
	for i := 0; i < b.N; i++ {
		f(10000, str)
	}
}

func BenchmarkPlusConcat(b *testing.B)       { benchmarkStrConcat(b, plusConcat) }
func BenchmarkSprintfConcat(b *testing.B)    { benchmarkStrConcat(b, sprintfConcat) }
func BenchmarkBuilderConcat(b *testing.B)    { benchmarkStrConcat(b, builderConcat) }
func BenchmarkBufferConcat(b *testing.B)     { benchmarkStrConcat(b, bufferConcat) }
func BenchmarkByteConcat(b *testing.B)       { benchmarkStrConcat(b, byteConcat) }
func BenchmarkPreByteConcat(b *testing.B)    { benchmarkStrConcat(b, preByteConcat) }
func BenchmarkPreBuilderConcat(b *testing.B) { benchmarkStrConcat(b, preBuilderConcat) }

// PreByteConcat 比 PlusConcat 快 1000+ 倍
//BenchmarkPlusConcat
//BenchmarkPlusConcat-8                 15          71833427 ns/op
//BenchmarkSprintfConcat
//BenchmarkSprintfConcat-8               8         126318812 ns/op
//BenchmarkBuilderConcat
//BenchmarkBuilderConcat-8            9400            114630 ns/op
//BenchmarkBufferConcat
//BenchmarkBufferConcat-8            10000            112820 ns/op
//BenchmarkByteConcat
//BenchmarkByteConcat-8               9225            122662 ns/op
//BenchmarkPreByteConcat
//BenchmarkPreByteConcat-8           21706             54849 ns/op
//BenchmarkPreBuilderConcat
//BenchmarkPreBuilderConcat-8        23762             50814 ns/op
