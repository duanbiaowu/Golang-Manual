package performance

import (
	"runtime"
	"testing"
	"time"
)

// 协程无法退出
func TestBadTimeout(t *testing.T) {
	t.Helper()
	for i := 0; i < 1000; i++ {
		_ = timeout(doBadThing)
	}
	time.Sleep(time.Second * 2)                         // 确保函数执行完毕
	t.Logf("NumGoroutine = %d", runtime.NumGoroutine()) // 1002
}

// 协程正常退出
func TestBufferTimeout(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = timeoutWithBuffer(doBadThing)
	}
	time.Sleep(time.Second * 2)
	t.Logf("NumGoroutine = %d", runtime.NumGoroutine()) // 2
}

func TestGoodTimeout(t *testing.T) {
	t.Helper()
	for i := 0; i < 1000; i++ {
		_ = timeout(doGoodThing)
	}
	time.Sleep(time.Second * 2)
	t.Logf("NumGoroutine = %d", runtime.NumGoroutine()) // 2
}

func TestPhasesTimeout(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_ = timeoutFirstPhase()
	}
	time.Sleep(time.Second * 3)
	t.Logf("NumGoroutine = %d", runtime.NumGoroutine()) // 2
}
