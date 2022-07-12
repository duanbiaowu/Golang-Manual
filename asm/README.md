# 快速入门

## 实现和声明
Go汇编语言并不是一个独立的语言，因为Go汇编程序无法独立使用。Go汇编代码必须以Go包的方式组织，同时包中至少要有一个Go语言文件用于
指明当前包名等基本包信息。如果Go汇编代码中定义的变量和函数要被其它Go语言代码引用，还需要通过Go语言代码将汇编中定义的符号声明出来。
用于变量的定义和函数的定义Go汇编文件类似于C语言中的.c文件，而用于导出汇编中定义符号的Go源文件类似于C语言的.h文件。


# reference
1. https://chai2010.cn/advanced-go-programming-book
2. https://developer.51cto.com/article/704916.html