package main

// CFLAGS通过-I./number将number库对应头文件所在的目录加入头文件检索路径
// LDFLAGS通过-L${SRCDIR}/number将编译后number静态库所在目录加为链接库检索路径，-lnumber表示链接libnumber.a静态库
// 在链接部分的检索路径不能使用相对路径（C/C++代码的链接程序所限制），
// 必须通过cgo特有的${SRCDIR}变量将源文件对应的当前目录路径展开为绝对路径（因此在windows平台中绝对路径不能有空白符号）

/*
#cgo CFLAGS: -I./number
#cgo LDFLAGS: -L${SRCDIR}/number -lnumber
#include "hello.h"
#include "number.h"
#include <stdio.h>

struct Person {
    int age;
    float height;
};

void SayName(const char* s);

static void SayHello(const char* s) {
    puts(s);
}

static void SayAge(int n) {
    printf("age: %d\n", n);
}

static void noreturn() {}
*/
import "C"
import (
	"fmt"
)

// 通过CGO的 //export SayHello 指令将Go语言实现的函数SayHello导出为C语言函数。为了适配CGO导出的C语言函数，
// 禁止了在函数的声明语句中的const修饰符。需要注意的是，这里其实有两个版本的SayHello函数：一个Go语言环境的；另一个是C语言环境的。
// cgo生成的C语言版本SayHello函数最终会通过桥接代码调用Go语言版本的SayHello函数。

// 执行的时候是先从Go语言的main函数，到CGO自动生成的C语言版本SayHello桥接函数，最后又回到了Go语言环境的SayHello函数

//export SayHello2
func SayHello2(s *C.char) {
	fmt.Println(C.GoString(s))
}

func main() {
	C.puts(C.CString("Hello, World"))
	C.SayHello(C.CString("Hello, World"))
	C.SayName(C.CString("Hello, Go"))
	C.SayNever(C.CString("Hello, Go"))
	C.SayHello2(C.CString("Hello, World"))

	C.SayAge(C.int(100))

	var person C.struct_Person
	fmt.Println(person.age)
	fmt.Println(person.height)

	res, _ := C.noreturn()
	fmt.Printf("%#v\n", res)

	fmt.Println(C.number_add_mod(10, 5, 12))
}
