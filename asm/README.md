# 快速入门

## 实现和声明
Go汇编语言并不是一个独立的语言，因为Go汇编程序无法独立使用。Go汇编代码必须以Go包的方式组织，同时包中至少要有一个Go语言文件用于
指明当前包名等基本包信息。如果Go汇编代码中定义的变量和函数要被其它Go语言代码引用，还需要通过Go语言代码将汇编中定义的符号声明出来。
用于变量的定义和函数的定义Go汇编文件类似于C语言中的.c文件，而用于导出汇编中定义符号的Go源文件类似于C语言的.h文件。

## 基本指令
### 栈调整
```shell
SUBQ $0x18, SP // 对 SP 做减法，为函数分配函数栈帧
ADDQ $0x18, SP // 对 SP 做加法，清除函数栈帧
```
### 数据搬运
搬运的长度是由 MOV 的后缀决定的
```shell
MOVB $1, DI      // 1 byte
MOVW $0x10, BX   // 2 bytes
MOVD $1, DX      // 4 bytes
MOVQ $-10, AX     // 8 bytes
```
### 常见计算指令
类似数据搬运指令，同样可以通过修改指令的后缀来对应不同长度的操作数。例如 ADDQ/ADDW/ADDL/ADDB。
```shell
ADDQ  AX, BX   // BX += AX
SUBQ  AX, BX   // BX -= AX
IMULQ AX, BX   // BX *= AX
```
### 条件跳转/无条件跳转 
```shell
// 无条件跳转
JMP addr   // 跳转到地址，地址可为代码中的地址，不过实际上手写不会出现这种东西
JMP label  // 跳转到标签，可以跳转到同一函数内的标签位置
JMP 2(PC)  // 以当前指令为基础，向前/后跳转 x 行
JMP -2(PC) // 同上

// 有条件跳转
JZ target // 如果 zero flag 被 set 过，则跳转
```

## 寄存器
### 通用寄存器
虽然 rbp 和 rsp 也可以用，不过 bp 和 sp 会被用来管理栈顶和栈底，最好不要拿来进行运算。
plan9 中使用寄存器不需要带 r 或 e 的前缀，例如 rax，只要写 AX 即可:
```shell
MOVQ $101, AX = mov rax, 101
```
### 伪寄存器
* FP: Frame pointer: arguments and locals.
* PC: Program counter: jumps and branches.
* SB: Static base pointer: global symbols.
* SP: Stack pointer: the highest address within the local stack frame.

补充说明:
FP: 使用形如 symbol+offset(FP) 的方式，引用函数的输入参数。例如 arg0+0(FP)，arg1+8(FP)，使用 FP 不加 symbol 时，无法通过编译，
在汇编层面来讲，symbol 并没有什么用，加 symbol 主要是为了提升代码可读性。另外，官方文档虽然将伪寄存器 FP 称之为 frame pointer，
实际上它根本不是 frame pointer，按照传统的 x86 的习惯来讲，frame pointer 是指向整个 stack frame 底部的 BP 寄存器。
假如当前的 callee 函数是 add，在 add 的代码中引用 FP，该 FP 指向的位置不在 callee 的 stack frame 之内，
而是在 caller 的 stack frame 上。具体可参见之后的 栈结构 一章。

PC: 实际上就是在体系结构的知识中常见的 pc 寄存器，在 x86 平台下对应 ip 寄存器，amd64 上则是 rip。除了个别跳转之外，
手写 plan9 代码与 PC 寄存器打交道的情况较少。

SB: 全局静态基指针，一般用来声明函数或全局变量，在之后的函数知识和示例部分会看到具体用法。
SP: plan9 的这个 SP 寄存器指向当前栈帧的局部变量的开始位置，使用形如 symbol+offset(SP) 的方式，引用函数的局部变量。
offset 的合法取值是 [-framesize, 0)，注意是个左闭右开的区间。假如局部变量都是 8 字节，那么第一个局部变量就可以用 localvar0-8(SP) 来表示。
这也是一个词不表意的寄存器。与硬件寄存器 SP 是两个不同的东西，在栈帧 size 为 0 的情况下，伪寄存器 SP 和硬件寄存器 SP 指向同一位置。
手写汇编代码时，如果是 symbol+offset(SP) 形式，则表示伪寄存器 SP。如果是 offset(SP) 则表示硬件寄存器 SP。务必注意。对于编译输出(go tool compile -S / go tool objdump)的代码来讲，目前所有的 SP 都是硬件寄存器 SP，无论是否带 symbol。

这里对容易混淆的几点简单进行说明：
1. 伪 SP 和硬件 SP 不是一回事，在手写代码时，伪 SP 和硬件 SP 的区分方法是看该 SP 前是否有 symbol。如果有 symbol，那么即为伪寄存器，
如果没有，那么说明是硬件 SP 寄存器。 
2. SP 和 FP 的相对位置是会变的，所以不应该尝试用伪 SP 寄存器去找那些用 FP + offset 来引用的值，例如函数的入参和返回值。 
3. 官方文档中说的伪 SP 指向 stack 的 top，是有问题的。其指向的局部变量位置实际上是整个栈的栈底(除 caller BP 之外)，
所以说 bottom 更合适一些。 
4. 在 go tool objdump/go tool compile -S 输出的代码中，是没有伪 SP 和 FP 寄存器的，我们上面说的区分伪 SP 和硬件 SP 寄存器的方法，
对于上述两个命令的输出结果是没法使用的。在编译和反汇编的结果中，只有真实的 SP 寄存器。 
5. FP 和 Go 的官方源代码里的 framepointer 不是一回事，源代码里的 framepointer 指的是 caller BP 寄存器的值，
在这里和 caller 的伪 SP 是值是相等的。

## 变量声明
```shell
# 使用 DATA 结合 GLOBL 来定义一个变量
# GLOBL 必须跟在 DATA 指令之后，使用 GLOBL 指令将变量声明为 global，额外接收两个参数，一个是 flag，另一个是变量的总大小。
# 大多数参数都是字面意思，不过这个 offset 需要稍微注意。其含义是该值相对于符号 symbol 的偏移，而不是相对于全局某个地址的偏移。
DATA  symbol+offset(SB)/width, value
GLOBL divtab(SB), RODATA, $64

# 所有符号在声明时，其 offset 一般都是 0
# 如果在全局变量中定义数组，或字符串，这时候就需要用上非 0 的 offset 了
# 新的标记 <>，这个跟在符号名之后，表示该全局变量只在当前文件中生效，类似于 C 语言中的 static。如果在另外文件中引用该变量的话，
# 会报 relocation target not found 错误。
DATA bio<>+0(SB)/8, $"oh yes i"
DATA bio<>+8(SB)/8, $"am here "
GLOBL bio<>(SB), RODATA, $16
```

## .s 和 .go 文件的全局变量互通
```shell
main.go
var version float32 = 1.0
func getVersion() float32

main.s
# ·version(SB)，表示该符号需要链接器来帮我们进行重定向(relocation)，如果找不到该符号，会输出 relocation target not found 的错误。
TEXT ·getVersion(SB),NOSPLIT,$0-4
    MOVQ ·version(SB), AX  
    MOVQ AX, ret+0(FP)
    RET
```

## 函数声明
```shell
# 为什么要叫 TEXT? 根据程序数据在文件中和内存中的分段，代码在二进制文件中存储在 .text 段中，这里也就是一种约定俗成的起名方式。
# 定义中的 pkgname 部分是可以省略的，非想写也可以写上。不过写上 pkgname 的话，在重命名 package 之后还需要改代码，所以推荐最好还是不要写。
// func add(a, b int) int
//   => 该声明定义在同一个 package 下的任意 .go 文件中
//   => 只有函数头，没有实现
TEXT pkgname·add(SB), NOSPLIT, $0-8
    MOVQ a+0(FP), AX
    MOVQ a+8(FP), BX
    ADDQ AX, BX
    MOVQ BX, ret+16(FP)
    RET

                              参数及返回值大小
                                  | 
 TEXT pkgname·add(SB),NOSPLIT,$32-32
       |        |               |
      包名     函数名         栈帧大小(局部变量+可能需要的额外调用函数的参数空间的总大小，但不包括调用其它函数时的 ret address 的大小)

```

## 地址运算
地址运算也是用 lea 指令，英文原意为 Load Effective Address，amd64 平台地址都是 8 个字节，所以直接就用 LEAQ 就好:
```shell
LEAQ (BX)(AX*8), CX
// 上面代码中的 8 代表 scale
// scale 只能是 0、2、4、8
// 如果写成其它值:
// LEAQ (BX)(AX*3), CX
// ./a.s:6: bad scale: 3

// 用 LEAQ 的话，即使是两个寄存器值直接相加，也必须提供 scale
// 下面这样是不行的
// LEAQ (BX)(AX), CX
// asm: asmidx: bad address 0/2064/2067
// 正确的写法是
LEAQ (BX)(AX*1), CX

// 在寄存器运算的基础上，可以加上额外的 offset
LEAQ 16(BX)(AX*1), CX

// 三个寄存器做运算，还是别想了
// LEAQ DX(BX)(AX*8), CX
// ./a.s:13: expected end of operand, found 
```

## 伪寄存器 SP 、伪寄存器 FP 和硬件寄存器 SP
伪 SP 和伪 FP 的相对位置是会变化的，手写时不应该用伪 SP 和 >0 的 offset 来引用数据，否则结果可能会出乎你的预料。

## global symbol: size 错误
```shell
NameData: missing Go type information for global symbol: size 8
```
错误提示汇编中定义的NameData符号没有类型信息。其实Go汇编语言中定义的数据并没有所谓的类型，每个符号只不过是对应一块内存而已，
因此NameData符号也是没有类型的。但是Go语言是再带垃圾回收器的语言，而Go汇编语言是工作在自动垃圾回收体系框架内的。
当Go语言的垃圾回收器在扫描到NameData变量的时候，无法知晓该变量内部是否包含指针，因此就出现了这种错误。错误的根本原因并不是NameData没有类型，
而是NameData变量没有标注是否会含有指针信息。
### 解决方案
1. 通过给NameData变量增加一个NOPTR标志，表示其中不会包含指针数据可以修复该错误
```shell
#include "textflag.h"

GLOBL ·NameData(SB),NOPTR,$8
```
2. 通过给·NameData变量在Go语言中增加一个不含指针并且大小为8个字节的类型来修改该错误：
```go
package pkg

var NameData [8]byte
var Name string
```

# reference
1. https://chai2010.cn/advanced-go-programming-book
2. https://developer.51cto.com/article/704916.html
3. https://go.xargin.com/docs/assembly/assembly