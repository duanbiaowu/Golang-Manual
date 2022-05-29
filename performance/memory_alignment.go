// See the doc: https://geektutu.com/post/hpg-struct-alignment.html

package performance

// 对齐保证(align guarantee)
// For a variable x of any type: unsafe.Alignof(x) is at least 1.
// For a variable x of struct type: unsafe.Alignof(x) is the largest of all the values unsafe.
// 	Alignof(x.f) for each field f of x, but at least 1.
// For a variable x of array type:
// 	unsafe.Alignof(x) is the same as the alignment of a variable of the array’s element type.

// A struct or array type has size zero if it contains no fields (or elements, respectively)
//	that have a size greater than zero.
// Two distinct zero-size variables may have the same address in memory.

type memoryArg struct {
	num1 int
	num2 int
}

type memoryArg2 struct {
	num1 int16
	num2 int32
}

// 对齐倍数 4，占用 8 bytes
type memoryArg3 struct {
	a int8  // 1 byte
	b int16 // 2 bytes
	c int32 // 4 bytes
}

// 对齐倍数 4，占用 12 bytes
type memoryArg4 struct {
	a int8  // 1 byte
	c int32 // 4 bytes
	b int16 // 2 bytes
}

// 空 struct{} 大小为 0，作为其他 struct 的字段时，一般不需要内存对齐。
// 但是有一种情况除外：即当 struct{} 作为结构体最后一个字段时，需要内存对齐。
// 因为如果有指针指向该字段, 返回的地址将在结构体之外，如果此指针一直存活不释放对应的内存，就会有内存泄露的问题（该内存不因结构体释放而释放）。

// 占用 8 bytes
type memoryArg5 struct {
	c int32
	a struct{}
}

// 占用 4 bytes
type memoryArg6 struct {
	a struct{}
	c int32
}
