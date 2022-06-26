package base

import (
	"bytes"
	"errors"
	"runtime"
	"strconv"
)

func GoroutineId() (int64, error) {
	var goroutineSpace = []byte("goroutine ")

	b := make([]byte, 128)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		return -1, errors.New("get current goroutine id failed")
	}
	return strconv.ParseInt(string(b[:i]), 10, 64)
}
