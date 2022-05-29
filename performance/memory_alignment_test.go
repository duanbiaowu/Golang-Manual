package performance

import (
	"testing"
	"unsafe"
)

func TestMemoryAlign(t *testing.T) {
	println(unsafe.Alignof(memoryArg{}))
	println(unsafe.Alignof(memoryArg2{}))
	println(unsafe.Alignof(memoryArg3{}))
	println(unsafe.Alignof(memoryArg4{}))

	println(unsafe.Sizeof(memoryArg3{}))
	println(unsafe.Sizeof(memoryArg4{}))

	println(unsafe.Sizeof(memoryArg5{}))
	println(unsafe.Sizeof(memoryArg6{}))
}
