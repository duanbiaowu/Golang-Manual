package singleflight

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/singleflight"
)

var (
	cnt int32
)

func getNumber() int {
	// 假设这里会对数据库进行调用, 模拟不同并发下耗时不同
	atomic.AddInt32(&cnt, 1)
	time.Sleep(time.Duration(cnt) * time.Millisecond)
	return 1024
}

func singleFlightGetNumber(sg *singleflight.Group) int {
	v, _, _ := sg.Do("getNumber", func() (interface{}, error) {
		return getNumber(), nil
	})
	return v.(int)
}

func TestTSingleFlight(t *testing.T) {
	var wg sync.WaitGroup

	now := time.Now()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			assert.Equal(t, 1024, getNumber())
		}()
	}
	wg.Wait()
	t.Logf("take times %s", time.Since(now))

	atomic.AddInt32(&cnt, -cnt)

	var sg singleflight.Group
	now = time.Now()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			assert.Equal(t, 1024, singleFlightGetNumber(&sg))
		}()
	}
	wg.Wait()
	t.Logf("take times %s", time.Since(now))
}
