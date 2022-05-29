package performance

import (
	"sync"
	"testing"
)

func benchmark(b *testing.B, x rw, read, write int) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for k := 0; k < read*100; k++ {
			wg.Add(1)
			go func() {
				x.read()
				wg.Done()
			}()
		}
		for k := 0; k < write*100; k++ {
			wg.Add(1)
			go func() {
				x.write()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

// 通过参数设定
// 读多写少(读占 90%)
// 读少写多(读占 10%)
// 读写一致(各占 50%)
func BenchmarkLockReadMore(b *testing.B)    { benchmark(b, &lock{}, 9, 1) }
func BenchmarkLockReadMoreRW(b *testing.B)  { benchmark(b, &rwLock{}, 9, 1) }
func BenchmarkLockWriteMore(b *testing.B)   { benchmark(b, &lock{}, 1, 9) }
func BenchmarkLockWriteMoreRW(b *testing.B) { benchmark(b, &rwLock{}, 1, 9) }
func BenchmarkLockEqual(b *testing.B)       { benchmark(b, &lock{}, 5, 5) }
func BenchmarkLockEqualRW(b *testing.B)     { benchmark(b, &rwLock{}, 5, 5) }

// 读写比为 9:1 时，读写锁的性能约为互斥锁的 8 倍
// 读写比为 1:9 时，读写锁性能相当
// 读写比为 5:5 时，读写锁的性能约为互斥锁的 3 倍
//BenchmarkLockReadMore
//BenchmarkLockReadMore-8               15         107940307 ns/op
//BenchmarkLockReadMoreRW
//BenchmarkLockReadMoreRW-8            128           9108039 ns/op
//BenchmarkLockWriteMore
//BenchmarkLockWriteMore-8              12          88025075 ns/op
//BenchmarkLockWriteMoreRW
//BenchmarkLockWriteMoreRW-8            16          88598662 ns/op
//BenchmarkLockEqual
//BenchmarkLockEqual-8                  13         114571238 ns/op
//BenchmarkLockEqualRW
//BenchmarkLockEqualRW-8                39          53688964 ns/op
