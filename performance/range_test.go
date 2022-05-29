package performance

import (
	"fmt"
	"testing"
)

// go test -v -bench='IntSlice$' -run='IntSlice$' .

func BenchmarkForIntSlice(b *testing.B) {
	nums := rangeInternalGenerateInt(1024 * 1024)
	for i := 0; i < b.N; i++ {
		n := len(nums)
		var tmp int
		for k := 0; k < n; k++ {
			tmp = nums[k]
		}
		_ = tmp
	}
}

func BenchmarkRangeIntSlice(b *testing.B) {
	nums := rangeInternalGenerateInt(1024 * 1024)
	for i := 0; i < b.N; i++ {
		var tmp int
		for _, num := range nums {
			tmp = num
		}
		_ = tmp
	}
}

func BenchmarkForStruct(b *testing.B) {
	var items [1024]RangeItem
	for i := 0; i < b.N; i++ {
		length := len(items)
		var tmp int
		for k := 0; k < length; k++ {
			tmp = items[k].id
		}
		_ = tmp
	}
}

func BenchmarkRangeIndexStruct(b *testing.B) {
	var items [1024]RangeItem
	for i := 0; i < b.N; i++ {
		var tmp int
		for k := range items {
			tmp = items[k].id
		}
		_ = tmp
	}
}

func BenchmarkRangeStruct(b *testing.B) {
	var items [1024]RangeItem
	for i := 0; i < b.N; i++ {
		var tmp int
		for _, item := range items {
			tmp = item.id
		}
		_ = tmp
	}
}

// for 的性能大约是 range (同时遍历下标和值) 的 400+ 倍
//BenchmarkForStruct
//BenchmarkForStruct-8             2149584               541.0 ns/op
//BenchmarkRangeIndexStruct
//BenchmarkRangeIndexStruct-8      2113270               553.0 ns/op
//BenchmarkRangeStruct
//BenchmarkRangeStruct-8              5083            227662 ns/op

// 用一个非常简单的例子来证明 range 迭代时，返回的是拷贝
func BenchmarkUseCopyRangeStruct(b *testing.B) {
	persons := []struct{ no int }{{no: 1}, {no: 2}, {no: 3}}
	for _, s := range persons {
		s.no += 10
	}
	for i := 0; i < len(persons); i++ {
		persons[i].no += 100
	}
	fmt.Println(persons) // [{101} {102} {103}]
}

// []*RangeItem
func BenchmarkForPointerStruct(b *testing.B) {
	items := rangeInternalGenerateItem(1024)
	for i := 0; i < b.N; i++ {
		length := len(items)
		var tmp int
		for k := 0; k < length; k++ {
			tmp = items[k].id
		}
		_ = tmp
	}
}

func BenchmarkRangePointerStruct(b *testing.B) {
	items := rangeInternalGenerateItem(1024)
	for i := 0; i < b.N; i++ {
		var tmp int
		for _, item := range items {
			tmp = item.id
		}
		_ = tmp
	}
}
