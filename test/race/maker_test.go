package race

import (
	"fmt"
	"testing"
)

type IceCreamMaker interface {
	// Great a customer.
	Hello()
}

type Ben struct {
	name string
}

func (b *Ben) Hello() {
	fmt.Printf("Ben says, \"Hello my name is %s\"\n", b.name)
}

type Jerry struct {
	name string
}

func (j *Jerry) Hello() {
	fmt.Printf("Jerry says, \"Hello my name is %s\"\n", j.name)
}

func TestMakerDataRace(t *testing.T) {
	var ben = &Ben{name: "Ben"}
	var jerry = &Jerry{"Jerry"}
	var maker IceCreamMaker = ben

	var loop0, loop1 func()

	loop0 = func() {
		maker = ben
		go loop1()
	}

	loop1 = func() {
		maker = jerry
		go loop0()
	}

	go loop0()

	for {
		maker.Hello()
	}
}

// 这个例子有趣的点在于，最后输出的结果会有这种例子:
// Ben says, "Hello my name is Jerry"
// Ben says, "Hello my name is Jerry"
// 因为我们在 maker = jerry  这种赋值操作的时候并不是原子的，只有对 single machine word 进行赋值的时候才是原子的，
// 虽然这个看上去只有一行，但是 interface 在 go 中其实是一个结构体，它包含了 type 和 data 两个部分，所以它的复制也不是原子的，会出现问题

// type interface struct {
//       Type uintptr     // points to the type of the interface implementation
//       Data uintptr     // holds the data for the interface's receiver
// }

// 这个案例有趣的点还在于，两个结构体的内存布局一模一样，所以出现错误也不会 panic 退出，
// 如果在里面再加入一个 string 的字段，去读取就会导致 panic，但是这也恰恰说明这个案例很可怕，这种错误在线上实在太难发现了，而且很有可能很致命。
