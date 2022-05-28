package pprof

import (
	"os"
	"runtime/pprof"
	"testing"
)

func TestBubbleSort(t *testing.T) {
	f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()

	err := pprof.StartCPUProfile(f)
	if err != nil {
		t.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	n := 10
	for i := 0; i < 5; i++ {
		nums := generate(n)
		bubbleSort(nums)
		n *= 10
	}
}
