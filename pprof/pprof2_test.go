package pprof

import (
	"net/http"
	_ "net/http/pprof"
	"testing"
)

func funcA() []byte {
	a := make([]byte, 10*1024*1024)
	return a
}

func funcB() ([]byte, []byte) {
	a := make([]byte, 10*1024*1024)
	b := funcA()
	return a, b
}

func funcC() ([]byte, []byte, []byte) {
	a := make([]byte, 10*1024*1024)
	b, c := funcB()
	return a, b, c
}

func TestPprof(t *testing.T) {
	for i := 0; i < 5; i++ {
		funcA()
		funcB()
		funcC()
	}

	_ = http.ListenAndServe("127.0.0.1:10001", nil)
}

// dump
// curl -sS 'http://127.0.0.1:10001/debug/pprof/heap?seconds=5' -o heap.pporf
