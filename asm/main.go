package main

import (
	"fmt"
	"unsafe"
)

type Text struct {
	Language string
	_        uint8
	Length   int
}

func add(a, b int) int
func addX(a, b int) int

// 获取 Text 的 Length 字段的值
func length(text *Text) int

// 获取 Text 结构体的大小
func sizeOfTextStruct() int

func main() {
	println(add(1, 2))
	println(addX(1, 2))
	text := &Text{
		Language: "Go",
		Length:   1024,
	}
	fmt.Println(text)
	println(length(text))
	println(sizeOfTextStruct())
	println(unsafe.Sizeof(*text))
}
