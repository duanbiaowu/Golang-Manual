# dlv with gdb
**推荐 dlv，因为 gdb 功能实在是有限，gdb 不理解 golang 的业务类型和协程。但是 gdb 有一个功能是无法替代的，就是 gcore 的功能。**

# dlv

## Install

```shell
go get github.com/go-delve/delve/cmd/dlv
```

## 示例
```shell
# 调试二进制
dlv exec <path/to/binary> [flags]
dlv exec ./example

# dlv 调试二进制，并带参数
dlv exec ./example -- --audit=./d

# 调试进程
dlv attach ${pid} [executable] [flags]

# 进程号是必选的。
dlv attach 12808 ./example

# 调试 core 文件
# dlv 调试core文件；并且标准输出导出到文件

dlv core <executable> <core> [flags]
dlv core ./example core.277282
```

## 调试常用语法
```shell
程序运行
call ：call 函数（注意了，这个会导致整个程序运行的）
continue ：往下运行
next ：单步调试
restart ：重启
step ：单步调试，某个函数
step-instruction ：单步调试某个汇编指令
stepout ：从当前函数跳出

断点相关
break (alias: b) ：设置断点
breakpoints (alias: bp)  ：打印所有的断点信息
clear ：清理断点
clearall ：清理所有的断点
condition (alias: cond)  ：设置条件断点
on ：设置一段命令，当断点命中的时候
trace (alias: t) ：设置一个跟踪点，这个跟踪点也是一个断点，只不过运行道德时候不会断住程序，只是打印一行信息，这个命令在某些场景是很有用的，比如你断住程序就会影响逻辑（业务有超时），而你仅仅是想打印某个变量而已，那么用这种类型的断点就行；；

信息打印
args : 打印程序的传参
examinemem (alias: x)  ：这个是神器，解析内存用的，和 gdb 的 x 命令一样；
locals ：打印本地变量
print (alias: p) ：打印一个表达式，或者变量
regs ：打印寄存器的信息
set ：set 赋值
vars ：打印全局变量（包变量）
whatis ：打印类型信息

协程相关
goroutine (alias: gr) ：打印某个特定协程的信息
goroutines (alias: grs)  ：列举所有的协程
thread (alias: tr) ：切换到某个线程
threads ：打印所有的线程信息

栈相关
deferred ：在 defer 函数上下文里执行命令
down ：上堆栈
frame ：跳到某个具体的堆栈
stack (alias: bt)  ：打印堆栈信息
up ：下堆栈

其他命令
config ：配置变更
disassemble (alias: disass) ：反汇编
edit (alias: ed) ：略
exit (alias: quit | q) ：略
funcs ：打印所有函数符号
libraries ：打印所有加载的动态库
list (alias: ls | l) ：显示源码
source ：加载命令
sources ：打印源码
types ：打印所有类型信息
```

# gdb
gdb 对 golang 的调试支持是通过一个 python 脚本文件 src/runtime/runtime-gdb.py 来扩展的，所以功能非常有限。
gdb 只能做到最基本的变量打印，却理解不了 golang 的一些特殊类型，比如 channel，map，slice 等，gdb 原生是无法调适 goroutine 协程的，
因为这个是用户态的调度单位，gdb 只能理解线程。所以只能通过 python 脚本的扩展，把协程结构按照链表输出出来，支持的命令：

**gdb当前只支持6个命令：**
```shell
3个 cmd 命令
info goroutines；打印所有的goroutines
goroutine ${id} bt；打印一个goroutine的堆栈
iface；打印静态或者动态的接口类型

3个函数
len；打印string，slices，map，channels 这四种类型的长度
cap；打印slices，channels 这两种类型的cap
dtype；强制转换接口到动态类型。
```

## 示例
```shell
# 打印全局变量 (注意单引号)
(gdb) p 'runtime.firstmoduledata'

# 打印数组变量长度
(gdb) p $len(xxx)
```

# 常见问题

## 调用上下文
```go
// 打印当前代码位置的堆栈
debug.PrintStack()
```

## 单点调试总是非预期的执行代码？
这种情况一般是被编译器优化了，比如函数内联了，编译出的二进制删减了无效逻辑、无效参数。这种情况就会导致你 dlv 单步调试的时候，总是非预期的执行，
或者打印某些变量打印不出来。这种情况解决方法就是：**禁止编译优化。**
```shell
go build -gcflags "-N -l"
```

# reference
1. https://mp.weixin.qq.com/s/OXpWRiCHcxpFlylUk3RDEA