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

TEXT ·length(SB),NOSPLIT,$0-16
    MOVQ text+0(FP), AX
    MOVQ Text_Length(AX), AX // 通过字段在结构体中的偏移量读取字段值
    MOVQ AX, ret+8(FP)
    RET

TEXT ·sizeOfTextStruct(SB),NOSPLIT,$0-8
    MOVQ $Text__size, AX // 保存结构体的大小到 AX 寄存器
    MOVQ AX, ret+0(FP)
    RET

GLOBL x(SB), RODATA, $8; // 声明全局变量 x
DATA  x+0(SB)/8, $10    // 初始化全局变量 x, 赋值为 10
