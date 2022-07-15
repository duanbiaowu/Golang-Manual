package main

import (
	"Golang-Patterns/asm/pkg"
	"fmt"
	"unsafe"
)

type Text struct {
	Language string
	_        uint8
	Length   int
}

var version float32 = 1.0

func add(a, b int) int
func addX(a, b int) int

func sub(a, b int) int
func mul(a, b int) int

// 获取 Text 的 Length 字段的值
func length(text *Text) int

// 获取 Text 结构体的大小
func sizeOfTextStruct() int

func getAge() int32
func getPI() float64
func getBirthYear() int32

func getVersion() float32

func main() {
	println(pkg.Id)
	println(pkg.Name)
	println(pkg.UserName)
	pkg.Say()

	println(add(1, 2))
	println(addX(1, 2))
	println(sub(10, 5))
	println(mul(10, 5))

	text := &Text{
		Language: "Go",
		Length:   1024,
	}
	fmt.Println(text)
	println(length(text))
	println(sizeOfTextStruct())
	println(unsafe.Sizeof(*text))

	println(getAge())
	println(getPI())
	println(getBirthYear())

	println(getVersion())
}
