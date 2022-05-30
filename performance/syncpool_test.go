package performance

import (
	"bytes"
	"encoding/json"
	"testing"
)

func BenchmarkStructUnmarshalNative(b *testing.B) {
	var buf []byte
	for n := 0; n < b.N; n++ {
		stu := &syncPoolStudent{}
		_ = json.Unmarshal(buf, stu)
	}
}

func BenchmarkStructUnmarshalWithPool(b *testing.B) {
	var buf []byte
	for n := 0; n < b.N; n++ {
		stu := studentPool.Get().(*syncPoolStudent)
		_ = json.Unmarshal(buf, stu)
		studentPool.Put(stu)
	}
}

// 因为 Student 结构体内存占用较小，内存分配几乎不耗时间。
// 而标准库 json 反序列化时利用了反射，效率是比较低的，占据了大部分时间。
// 内存占用差了一个数量级，使用了 sync.Pool 后，内存占用仅为未使用的 168/1320 = 1/12，对 GC 的影响就很大。
//BenchmarkStructUnmarshalNative
//BenchmarkStructUnmarshalNative-8         3258724               376.9 ns/op          1320 B/op          3 allocs/op
//BenchmarkStructUnmarshalWithPool
//BenchmarkStructUnmarshalWithPool-8      10093932               122.3 ns/op           168 B/op          2 allocs/op

func BenchmarkByteBuffer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var buf bytes.Buffer
		buf.Write(data)
	}
}

func BenchmarkByteBufferWithPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Write(data)
		buf.Reset()
		bufferPool.Put(buf)
	}
}

//BenchmarkByteBuffer
//BenchmarkByteBuffer-8                     697396              1742 ns/op           10240 B/op          1 allocs/op
//BenchmarkByteBufferWithPool
//BenchmarkByteBufferWithPool-8            9740994               108.2 ns/op             0 B/op          0 allocs/op
