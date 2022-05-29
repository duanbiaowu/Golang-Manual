// See the doc: https://geektutu.com/post/hpg-mutex.html

package performance

import (
	"sync"
	"time"
)

type rw interface {
	write()
	read()
}

type lock struct {
	count int
	mu    sync.Mutex
}

type rwLock struct {
	count int
	mu    sync.RWMutex
}

const cost = time.Microsecond

func (l *lock) write() {
	l.mu.Lock()
	l.count++
	time.Sleep(cost)
	l.mu.Unlock()
}

func (l *lock) read() {
	l.mu.Lock()
	_ = l.count
	time.Sleep(cost)
	l.mu.Unlock()
}

func (l *rwLock) write() {
	l.mu.Lock()
	l.count++
	time.Sleep(cost)
	l.mu.Unlock()
}

func (l *rwLock) read() {
	l.mu.RLock()
	_ = l.count
	time.Sleep(cost)
	l.mu.RUnlock()
}
