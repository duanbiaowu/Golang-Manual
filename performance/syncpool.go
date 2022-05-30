// See the doc: https://geektutu.com/post/hpg-sync-pool.html

package performance

import (
	"bytes"
	"sync"
)

type syncPoolStudent struct {
	Name   string
	Age    int32
	Remark [1024]byte
}

var (
	studentPool = sync.Pool{
		New: func() interface{} {
			return new(syncPoolStudent)
		},
	}

	bufferPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}

	data = make([]byte, 10240)
)
