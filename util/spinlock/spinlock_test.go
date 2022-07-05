// Copyright 2021 Andy Pan & Dietoad. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package spinlock

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// goos: linux
// goarch: amd64
// pkg: Golang-Patterns/util/spinlock
// cpu: Intel(R) Core(TM) i5-8300H CPU @ 2.30GHz
// BenchmarkMutex
// BenchmarkMutex-8                21886387                55.83 ns/op
// BenchmarkSpinLock
// BenchmarkSpinLock-8             46848830                25.81 ns/op
// BenchmarkBackOffSpinLock
// BenchmarkBackOffSpinLock-8      55894545                21.16 ns/op

type originSpinLock uint32

func (sl *originSpinLock) Lock() {
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		runtime.Gosched()
	}
}

func (sl *originSpinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}

func NewOriginSpinLock() sync.Locker {
	return new(originSpinLock)
}

func BenchmarkMutex(b *testing.B) {
	m := sync.Mutex{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Lock()
			//nolint:static check
			m.Unlock()
		}
	})
}

func BenchmarkSpinLock(b *testing.B) {
	spin := NewOriginSpinLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			//nolint:static check
			spin.Unlock()
		}
	})
}

func BenchmarkBackOffSpinLock(b *testing.B) {
	spin := NewSpinLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			//nolint:static check
			spin.Unlock()
		}
	})
}
