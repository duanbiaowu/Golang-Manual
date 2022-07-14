#include "textflag.h"
#include "go_asm.h" // 该文件自动生成

TEXT ·add(SB),NOSPLIT,$0-24
    MOVQ a+0(FP), AX // 读取第一个参数
    MOVQ b+8(FP), BX // 读取第二个参数
    ADDQ BX, AX
    MOVQ AX, ret+16(FP) // 保存结果
    RET

TEXT ·addX(SB),NOSPLIT,$0-24
    MOVQ a+0(FP), AX
    MOVQ b+8(FP), BX
    ADDQ BX, AX
    MOVQ $x(SB), BX // 读取全局变量 x 的地址
    MOVQ 0(BX), BX  // 读取全局变量 x 的值
    ADDQ BX, AX
    MOVQ AX, ret+16(FP)
    RET

TEXT ·sub(SB), NOSPLIT, $0-24
    MOVQ a+0(FP), AX
    MOVQ b+8(FP), BX
    SUBQ BX, AX    // AX -= BX
    MOVQ AX, ret+16(FP)
    RET

// func mul(a, b int) int
TEXT ·mul(SB), NOSPLIT, $0-24
    MOVQ  a+0(FP), AX
    MOVQ  b+8(FP), BX
    IMULQ BX, AX    // AX *= BX
    MOVQ  AX, ret+16(FP)
    RET

TEXT ·length(SB),NOSPLIT,$0-16
    MOVQ text+0(FP), AX
    MOVQ Text_Length(AX), AX // 通过字段在结构体中的偏移量读取字段值
    MOVQ AX, ret+8(FP)
    RET

TEXT ·sizeOfTextStruct(SB),NOSPLIT,$0-8
    MOVQ $Text__size, AX // 保存结构体的大小到 AX 寄存器
    MOVQ AX, ret+0(FP)
    RET

TEXT ·getAge(SB),NOSPLIT,$0-4
    MOVQ age(SB), AX
    MOVQ AX, ret+0(FP)
    RET

TEXT ·getPI(SB),NOSPLIT,$0-8
    MOVQ pi(SB), AX
    MOVQ AX, ret+0(FP)
    RET

TEXT ·getBirthYear(SB),NOSPLIT,$0-4
    MOVQ birthYear(SB), AX
    MOVQ AX, ret+0(FP)
    RET

TEXT ·getVersion(SB),NOSPLIT,$0-4
    MOVQ ·version(SB), AX
    MOVQ AX, ret+0(FP)
    RET

DATA  x+0(SB)/8, $10    // 初始化全局变量 x, 赋值为 10
GLOBL x(SB), RODATA, $8 // 声明全局变量 x, GLOBL 必须跟在 DATA 指令之后

DATA age+0x00(SB)/4, $18
GLOBL age(SB), RODATA, $4

DATA pi+0(SB)/8, $3.1415926
GLOBL pi(SB), RODATA, $8

DATA birthYear+0(SB)/4, $1992
GLOBL birthYear(SB), RODATA, $4

// 最后一行的空行是必须的，否则可能报 unexpected EOF，或者在最后一行代码加 ;