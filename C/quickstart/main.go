package main

/*
#include "hello.h"
#include <stdio.h>

void SayName(const char* s);

static void SayHello(const char* s) {
    puts(s);
}
*/
import "C"
import "fmt"

// 通过CGO的 //export SayHello 指令将Go语言实现的函数SayHello导出为C语言函数。为了适配CGO导出的C语言函数，
// 禁止了在函数的声明语句中的const修饰符。需要注意的是，这里其实有两个版本的SayHello函数：一个Go语言环境的；另一个是C语言环境的。
// cgo生成的C语言版本SayHello函数最终会通过桥接代码调用Go语言版本的SayHello函数。

// 执行的时候是先从Go语言的main函数，到CGO自动生成的C语言版本SayHello桥接函数，最后又回到了Go语言环境的SayHello函数

//export SayHello2
func SayHello2(s *C.char) {
	fmt.Print(C.GoString(s))
}

func main() {
	C.puts(C.CString("Hello, World"))
	C.SayHello(C.CString("Hello, World"))
	C.SayName(C.CString("Hello, Go"))
	C.SayNever(C.CString("Hello, Go"))
	C.SayHello2(C.CString("Hello, World"))
}
