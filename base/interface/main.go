package main

import "fmt"

type Country struct {
	Name string
}

type City struct {
	Name string
}

type StringAble interface {
	ToString() string
}

func (c Country) ToString() string {
	return "Country = " + c.Name
}

func (c City) ToString() string {
	return "City = " + c.Name
}

func PrintStr(p StringAble) {
	fmt.Println(p.ToString())
}

type Shape interface {
	Sides() int
	Area() int
}

type Square struct {
	len int
}

func (s *Square) Sides() int {
	return s.len * s.len
}

func main() {
	d1 := Country{"USA"}
	d2 := City{"Los Angeles"}
	PrintStr(d1)
	PrintStr(d2)

	// 接口完整性检查
	// 可以看到 Square 并没有实现 Shape 接口的所有方法
	// 程序虽然可以跑通，但是这样编程的方式并不严谨
	s := Square{len: 5}
	fmt.Printf("%d\n", s.Sides())

	// 强制实现接口的所有方法
	// 1. 声明一个 _ 变量 (不使用)
	// 2. 把一个 nil 指针从 Square 转为 Shape
	// 3. 如果 Square 没有实现 Shape 的全部方法，编译器就会报错:
	// Cannot use '(*Square)(nil)' (type *Square)
	// as the type Shape Type does not implement 'Shape' as some methods
	// are missing: Area() int
	//var _ Shape = (*Square)(nil)
}
