package performance

import (
	"runtime"
	"testing"
)

func printMem(t *testing.T) {
	t.Helper()
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	t.Logf("%.2f MB", float64(rtm.Alloc)/1024./1024.)
}

func testLastChars(t *testing.T, f func([]int) []int) {
	t.Helper()
	ans := make([][]int, 0)
	for k := 0; k < 100; k++ {
		// 64位，一个 int 占 8 Byte，128 * 1024 个整数恰好占据 1 MB
		origin := generateWithCap(128 * 1024)
		ans = append(ans, f(origin))
	}
	printMem(t)
}

func testLastCharsWithCallGC(t *testing.T, f func([]int) []int) {
	t.Helper()
	ans := make([][]int, 0)
	for k := 0; k < 100; k++ {
		// 64位，一个 int 占 8 Byte，128 * 1024 个整数恰好占据 1 MB
		origin := generateWithCap(128 * 1024)
		ans = append(ans, f(origin))
		runtime.GC()
	}
	printMem(t)
}

func TestLastCharsBySlice(t *testing.T)           { testLastChars(t, lastNumsBySlice) }
func TestLastCharsByCopy(t *testing.T)            { testLastChars(t, lastNumsByCopy) }
func TestLastCharsWithCallGCBySlice(t *testing.T) { testLastCharsWithCallGC(t, lastNumsBySlice) }
func TestLastCharsWithCallGCByCopy(t *testing.T)  { testLastCharsWithCallGC(t, lastNumsByCopy) }
