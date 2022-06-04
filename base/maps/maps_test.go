// See the doc: https://shaffer.cn/golang/golang-map-benchmark/
package maps

import (
	"strconv"
	"sync"
	"testing"

	cmap "github.com/orcaman/concurrent-map"
)

var N = 10000

var simpleMap map[string]interface{}
var mutexMap MutexMap
var syncMap sync.Map
var concurrentMap cmap.ConcurrentMap
var concurrentMap2 ConcurrentMap
var concurrentMap4 ConcurrentMap

type MutexMap struct {
	sync.RWMutex
	m map[string]interface{}
}

func (m *MutexMap) Set(k string, v interface{}) {
	m.Lock()
	m.m[k] = v
	m.Unlock()
}

func (m *MutexMap) Get(k string) (v interface{}) {
	m.RLock()
	v = m.m[k]
	m.RUnlock()
	return
}

// go test -bench=. -benchtime=3s -benchmem
func init() {
	simpleMap = make(map[string]interface{}, N)
	mutexMap = MutexMap{m: make(map[string]interface{}, N)}
	syncMap = sync.Map{}
	for i := 0; i < N; i++ {
		syncMap.Store(strconv.Itoa(i), i)
	}
	concurrentMap = cmap.New()
	for i := 0; i < N; i++ {
		concurrentMap.Set(strconv.Itoa(i), i)
	}
	concurrentMap2 = New(1)
	for i := 0; i < N; i++ {
		concurrentMap2.Set(strconv.Itoa(i), i)
	}
	concurrentMap4 = New(2)
	for i := 0; i < N; i++ {
		concurrentMap4.Set(strconv.Itoa(i), i)
	}
}

func BenchmarkMap_Set(b *testing.B) {
	finished := make(chan bool, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			simpleMap[strconv.Itoa((i+j)%N)] = i
		}
		finished <- true
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMap_Get(b *testing.B) {
	finished := make(chan bool, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_ = simpleMap[strconv.Itoa((i+j)%N)]
		}
		finished <- true
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMutexMap_Set(b *testing.B) {
	finished := make(chan bool, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				mutexMap.Set(strconv.Itoa((i+j)%N), i)
			}
			finished <- true
		}(i)
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMutexMap_Get(b *testing.B) {
	finished := make(chan bool, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				mutexMap.Get(strconv.Itoa((i + j) % N))
			}
			finished <- true
		}(i)
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkSyncMap_Set(b *testing.B) {
	finished := make(chan bool, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				syncMap.Store(strconv.Itoa((i+j)%N), i)
			}
			finished <- true
		}(i)
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkSyncMap_Get(b *testing.B) {
	finished := make(chan bool, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				syncMap.Load(strconv.Itoa((i + j) % N))
			}
			finished <- true
		}(i)
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkConcurrentMap_Set(b *testing.B) {
	finished := make(chan bool, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				concurrentMap.Set(strconv.Itoa((i+j)%N), i)
			}
			finished <- true
		}(i)
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkConcurrentMap_Get(b *testing.B) {
	finished := make(chan bool, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				concurrentMap.Get(strconv.Itoa((i + j) % N))
			}
			finished <- true
		}(i)
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkConcurrentMap2_Set(b *testing.B) {
	finished := make(chan bool, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				concurrentMap2.Set(strconv.Itoa((i+j)%N), i)
			}
			finished <- true
		}(i)
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkConcurrentMap2_Get(b *testing.B) {
	finished := make(chan bool, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				concurrentMap2.Get(strconv.Itoa((i + j) % N))
			}
			finished <- true
		}(i)
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkConcurrentMap4_Set(b *testing.B) {
	finished := make(chan bool, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				concurrentMap4.Set(strconv.Itoa((i+j)%N), i)
			}
			finished <- true
		}(i)
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkConcurrentMap4_Get(b *testing.B) {
	finished := make(chan bool, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				concurrentMap4.Get(strconv.Itoa((i + j) % N))
			}
			finished <- true
		}(i)
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}
