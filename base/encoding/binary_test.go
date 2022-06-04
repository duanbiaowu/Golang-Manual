// See the doc: https://zhuanlan.zhihu.com/p/35326716

package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// 基于文本类型的协议（比如 JSON）和二进制协议都是字节通信，他们不同点在于他们使用哪种类型的字节和如何组织这些字节。
// 文本协议只适用于 ASCII 或 Unicode 编码可打印的字符通信。
// 	例如 "26" 使用 "2" 和 "6" 的 utf 编码的字符串表示，这种方式方便我们读，但对于计算机效率较低。
// 在二进制协议中，同样数字 "26" 可使用一个字节 0x1A 十六进制表示，减少了一半的存储空间且原始的字节格式能够被计算机直接识别而不需解析。
//	当一个数字足够大的时候，性能优势就会明显体现。

// Big Endian 是指低地址端 存放 高位字节。
// Little Endian 是指低地址端 存放 低位字节。
// 以一个 32 位的 int 型变量 a = 0x12345678 举例说明
// Big Endian
// -----------------------------
// | 0x12 | 0x34 | 0x56 | 0x78 |
// -----------------------------
// 低地址 -------------> 高地址
// Litter Endian
// -----------------------------
// | 0x78 | 0x56 | 0x34 | 0x12 |
// -----------------------------
// 低地址 -------------> 高地址

// 各自优势：
// Big Endian：符号位的判定固定为第一个字节，容易判断正负。
// Little Endian：长度为1，2，4字节的数，排列方式都是一样的，数据类型转换非常方便。
// 在计算机内部，小端序被广泛应用于现代性 CPU 内部存储数据；而在其他场景譬如网络传输和文件存储使用大端序。

// 转化成二进制格式与原本数据转字符串相比会更节省空间
// binary 包实现了对数据与 byte 之间的转换，以及 varint 的编解码

func TestRead(t *testing.T) {
	var pi float64
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &pi)
	assert.Nil(t, err)
	assert.Equal(t, 3.141592653589793, pi)

	var data struct {
		PI   float64
		Uate uint8
		Mine [3]byte
		Too  uint16
	}

	b2 := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40, 0xff, 0x01, 0x02, 0x03, 0xbe, 0xef}
	err = binary.Read(bytes.NewReader(b2), binary.LittleEndian, &data)
	assert.Nil(t, err)

	assert.Equal(t, 3.141592653589793, data.PI)
	assert.Equal(t, uint8(255), data.Uate)
	assert.Equal(t, [3]uint8{1, 2, 3}, data.Mine)
	assert.Equal(t, uint16(61374), data.Too)
}

func TestWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	var pi float64 = math.Pi
	err := binary.Write(buf, binary.LittleEndian, pi)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}, buf.Bytes())

	buf2 := new(bytes.Buffer)
	var data = []interface{}{
		uint16(61374),
		int8(-54),
		uint8(254),
	}
	for _, v := range data {
		err = binary.Write(buf2, binary.LittleEndian, v)
		assert.Nil(t, err)
	}
	assert.Equal(t, []byte{0xbe, 0xef, 0xca, 0xfe}, buf2.Bytes())
}

func TestSize(t *testing.T) {
	var a int
	p := &a
	b := [10]int64{1}
	s := "adsa"
	bs := make([]byte, 10)

	assert.Equal(t, -1, binary.Size(a))
	assert.Equal(t, -1, binary.Size(p))
	assert.Equal(t, 80, binary.Size(b))
	assert.Equal(t, -1, binary.Size(s))
	assert.Equal(t, 10, binary.Size(bs))
}

// 将 uint64 类型放入 buf 中，并返回写入的字节数
// 如果buf过小，PutUvarint 抛出 panic
func TestPutVarint(t *testing.T) {
	buf := make([]byte, binary.MaxVarintLen64)
	for _, x := range []int64{-65, -64, -2, -1, 0, 1, 2, 63, 64} {
		n := binary.PutVarint(buf, x)
		fmt.Printf("src value bytes = %d, dest value bytes = %d\n", unsafe.Sizeof(x), n)
	}
}

// Uvarint 从 buf 中解码并返回一个uint64的数据,及解码的字节数(>0)
// 如果出错,则返回数据 0 和一个小于等于 0 的字节数 n
// 	1> n == 0: buf 太小
//	2> n < 0: 数据太大,超出 uint64 最大范围,且 -n 为已解析字节数
func TestUvarint(t *testing.T) {
	inputs := [][]byte{
		{0x01},
		{0x02},
		{0x7f},
		{0x80, 0x01},
		{0xff, 0x01},
		{0x80, 0x02},
	}
	for _, b := range inputs {
		_, n := binary.Uvarint(b)
		assert.Equal(t, len(b), n)
	}
}

// ReadUvarint 从 r 中解析并返回一个 uint64 类型的数据及出现的错误
func TestReadUvarint(t *testing.T) {
	var buf []byte
	var buf2 = []byte{144, 192, 192, 129, 132, 136, 140, 144, 16, 0, 1, 1}

	n, err := binary.ReadUvarint(bytes.NewBuffer(buf))
	assert.Equal(t, 0, n)
	assert.Equal(t, err, io.EOF)

	n, err = binary.ReadUvarint(bytes.NewBuffer(buf2))
	fmt.Println(n, err)
	assert.Equal(t, uint64(1161981756374523920), n)
	assert.Nil(t, err)
}
