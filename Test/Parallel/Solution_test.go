package golang

import (
	"bytes"
	"testing"
)

var pairs = []struct {
	k string
	v string
}{
	{"polaris", " 徐新华 "},
	{"studygolang", "Go 语言中文网 "},
	{"stdlib", "Go 语言标准库 "},
	{"polaris1", " 徐新华 1"},
	{"studygolang1", "Go 语言中文网 1"},
	{"stdlib1", "Go 语言标准库 1"},
	{"polaris2", " 徐新华 2"},
	{"studygolang2", "Go 语言中文网 2"},
	{"stdlib2", "Go 语言标准库 2"},
	{"polaris3", " 徐新华 3"},
	{"studygolang3", "Go 语言中文网 3"},
	{"stdlib3", "Go 语言标准库 3"},
	{"polaris4", " 徐新华 4"},
	{"studygolang4", "Go 语言中文网 4"},
	{"stdlib4", "Go 语言标准库 4"},
}

// 注释掉 WriteToMap 和 ReadFromMap 中 locker 保护的代码，
// 同时注释掉测试代码中的 t.Parallel，执行测试，测试通过，即使加上 -race，测试依然通过

// 只注释掉 WriteToMap 和 ReadFromMap 中 locker 保护的代码，
// 执行测试，测试失败（如果未失败，加上 -race 一定会失败）

func TestWriteToMap(t *testing.T) {
	t.Parallel()
	for _, tt := range pairs {
		WriteToMap(tt.k, tt.v)
	}
}

func TestReadFromMap(t *testing.T) {
	t.Parallel()
	for _, tt := range pairs {
		actual := ReadFromMap(tt.k)
		if actual != tt.v {
			t.Errorf("the value of key(%s) is %s, expected: %s", tt.k, actual, tt.v)
		}
	}
}

func BenchmarkWriteToMap(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine has its own bytes.Buffer.
		var buf bytes.Buffer
		for pb.Next() {
			// The loop body is executed b.N times total across all goroutines.
			buf.Reset()
		}
	})
}
