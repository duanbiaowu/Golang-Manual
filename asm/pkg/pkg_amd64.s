#include "textflag.h"

GLOBL ·Id(SB),NOPTR,$8
DATA ·Id+0(SB)/1,$0x37  // 高 1 位
DATA ·Id+1(SB)/1,$0x25  // ...
DATA ·Id+2(SB)/1,$0x00
DATA ·Id+3(SB)/1,$0x00
DATA ·Id+4(SB)/1,$0x00
DATA ·Id+5(SB)/1,$0x00
DATA ·Id+6(SB)/1,$0x00  // ...
DATA ·Id+7(SB)/1,$0x00  // 低 1 位

GLOBL ·NameData(SB),NOPTR,$8    // 通过给·NameData增加NOPTR标志的方式表示其中不含指针数据
DATA  ·NameData+0(SB)/8,$"gopher"

GLOBL ·Name(SB),NOPTR,$16
DATA  ·Name+0(SB)/8,$·NameData(SB)
DATA  ·Name+8(SB)/8,$6

// 在用汇编定义字符串时我们可以换一种思维：将底层的 字符串数据 和 字符串头结构体 定义在一起，这样可以避免引入NameData符号：
// 在新的结构中，UserName 符号对应的内存从16字节变为24字节，多出的8个字节存放底层的“gopher”字符串。
// ·UserName 符号前16个字节依然对应reflect.StringHeader结构体：Data部分对应 $·UserName+16(SB)，
// 表示数据的地址为 UserName 符号往后偏移16个字节的位置；Len部分依然对应6个字节的长度。这是C语言程序员经常使用的技巧。
GLOBL ·UserName(SB),NOPTR,$24
DATA ·UserName+0(SB)/8,$·UserName+16(SB)
DATA ·UserName+8(SB)/8,$6
DATA ·UserName+16(SB)/8,$"gopher"

TEXT ·Say(SB), NOSPLIT, $16-0
    MOVQ ·helloWorld+0(SB), AX
    MOVQ AX, 0(SP)
    MOVQ ·helloWorld+8(SB), BX
    MOVQ BX, 8(SP)
    CALL ·output(SB)    // 在调用 output 之前，已经把参数都通过物理寄存器 SP 搬到了函数的栈顶
    RET
