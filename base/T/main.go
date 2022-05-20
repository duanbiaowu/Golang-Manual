//go:build ignore
// +build ignore

package main

import (
	"fmt"
)

// 编译加上参数 -gcflags=-G=3
//（这个编译参数会在1.18版上成为默认参数）
// 泛型函数不支持 export，函数名第一个字符只能为小写

// Go 泛型基本可用，但是还有 3 个问题：
// 1. fmt.Printf() 中的泛型类型是 %v 还不够好，不能像c++ iostream 重载 >> 来获得程序自定义的输出
// 2. 不支持操作符重载，很难在泛型算法中使用 {泛型操作符} 如：==
// 3. 算法依赖于具体的数据结构，对于不同数据结构要重写。 没有一个像 C++ STL 的一个泛型迭代器。

// 1. Simple Example ---------------------------------------------
func print[T any] (arr []T) {
	for _, v := range arr {
		fmt.Print(v)
		fmt.Print(" ")
	}
	fmt.Println("")
}

func find[T comparable] (arr []T, elem T) int {
	for i, v := range arr {
		if v == elem {
			return i
		}
	}
	return -1
}

//strs := []string{"Hello", "World",  "Generics"}
//decs := []float64{3.14, 1.14, 1.618, 2.718 }
//nums := []int{2,4,6,8}
//print(strs)
//print(decs)
//print(nums)

// 2. Stack ---------------------------------------------
type stack [T any] []T

func (s *stack[T]) push(elem T) {
	*s = append(*s, elem)
}

func (s *stack[T]) pop() {
	if len(*s) > 0 {
		*s = (*s)[:len(*s)-1]
	}
}

// 如果栈为空，top() 返回 error | nil
// 除了类型泛型外，还需要有一些值泛型
// 所以这里返回一个指针，可以判断一下指针是否为空
func (s *stack[T]) top() *T {
	if len(*s) > 0 {
		return &(*s)[len(*s)-1]
	}
	return nil
}

func (s *stack[T]) len() int {
	return len(*s)
}

func (s *stack[T]) print() {
	for _, elem := range *s {
		fmt.Print(elem)
		fmt.Print(" ")
	}
	fmt.Println("")
}

//ss := stack[string]{}
//ss.push("Hello")
//ss.push("Hao")
//ss.push("Chen")
//ss.print()
//fmt.Printf("stack top is - %v\n", *(ss.top()))
//ss.pop()
//ss.pop()
//ss.print()
//
//ns := stack[int]{}
//ns.push(10)
//ns.push(20)
//ns.print()
//ns.pop()
//ns.print()
//*ns.top() += 1
//ns.print()
//ns.pop()
//fmt.Printf("stack top is - %v\n", ns.top())

// 3. LinkList ---------------------------------------------
type node[T comparable] struct {
	data T
	prev *node[T]
	next *node[T]
}

type list[T comparable] struct {
	head, tail *node[T]
	len        int
}

func (l *list[T]) isEmpty() bool {
	return l.head == nil && l.tail == nil
}

func (l *list[T]) add(data T) {
	n := &node[T]{
		data: data,
		prev: nil,
		next: l.head,
	}
	if l.isEmpty() {
		l.head = n
		l.tail = n
	}
	l.head.prev = n
	l.head = n
}

func (l *list[T]) push(data T) {
	n := &node[T]{
		data: data,
		prev: l.tail,
		next: nil,
	}
	if l.isEmpty() {
		l.head = n
		l.tail = n
	}
	l.tail.next = n
	l.tail = n
}

func (l *list[T]) del(data T) {
	for p := l.head; p != nil; p = p.next {
		if data == p.data {
			if p == l.head {
				l.head = p.next
			}
			if p == l.tail {
				l.tail = p.prev
			}
			if p.prev != nil {
				p.prev.next = p.next
			}
			if p.next != nil {
				p.next.prev = p.prev
			}
			return
		}
	}
}

func (l *list[T]) print() {
	if l.isEmpty() {
		fmt.Println("the link list is empty.")
		return
	}
	for p := l.head; p != nil; p = p.next {
		fmt.Printf("[%v] -> ", p.data)
	}
	fmt.Println("nil")
}

//var l = list[int]{}
//l.add(1)
//l.add(2)
//l.push(3)
//l.push(4)
//l.add(5)
//l.print() //[5] -> [2] -> [1] -> [3] -> [4] -> nil
//l.del(5)
//l.del(1)
//l.del(4)
//l.print() //[2] -> [3] -> nil

// 4. 泛型 Map ---------------------------------------------
func gMap[T1 any, T2 any] (arr []T1, f func(T1) T2) []T2 {
	result := make([]T2, len(arr))
	for i, elem := range arr {
		result[i] = f(elem)
	}
	return result
}

//nums := []int {0,1,2,3,4,5,6,7,8,9}
//squares := gMap(nums, func (elem int) int {
//	return elem * elem
//})
//print(squares)  //0 1 4 9 16 25 36 49 64 81
//
//strs := []string{"Hao", "Chen", "MegaEase"}
//upstrs := gMap(strs, func(s string) string  {
//	return strings.ToUpper(s)
//})
//print(upstrs) // HAO CHEN MEGAEASE
//
//
//dict := []string{"零", "壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"}
//strs =  gMap(nums, func (elem int) string  {
//	return  dict[elem]
//})
//print(strs) // 零 壹 贰 叁 肆 伍 陆 柒 捌 玖

// 5. 泛型 Reduce ---------------------------------------------
func gReduce[T1 any, T2 any] (arr []T1, init T2, f func(T2, T1) T2) T2 {
	result := init
	for _, elem := range arr {
		result = f(result, elem)
	}
	return result
}

//nums := []int {0,1,2,3,4,5,6,7,8,9}
//sum := gReduce(nums, 0, func (result, elem int) int  {
//	return result + elem
//})
//fmt.Printf("Sum = %d \n", sum)

// 6. 泛型 Filter ---------------------------------------------
func gFilter[T any] (arr []T, in bool, f func(T) bool) []T {
	result := []T{}
	for _, elem := range arr {
		choose := f(elem)
		if (in && choose) || (!in && !choose) {
			result = append(result, elem)
		}
	}
	return result
}

func gFilterIn[T any] (arr []T, f func(T) bool) []T {
	return gFilter(arr, true, f)
}

func gFilterOut[T any] (arr []T, f func(T) bool) []T {
	return gFilter(arr, false, f)
}

//nums := []int {0,1,2,3,4,5,6,7,8,9}
//odds := gFilterIn(nums, func (elem int) bool  {
//	return elem % 2 == 1
//})
//print(odds)

// 6. 泛型业务示例 ---------------------------------------------
type Employee struct {
	Name     string
	Age      int
	Vacation int
	Salary   float32
}

var employees = []Employee{
	{"Hao", 44, 0, 8000.5},
	{"Bob", 34, 10, 5000.5},
	{"Alice", 23, 5, 9000.0},
	{"Jack", 26, 0, 4000.0},
	{"Tom", 48, 9, 7500.75},
	{"Marry", 29, 0, 6000.0},
	{"Mike", 32, 8, 4000.3},
}

// 统计所有员工薪水
// gReduce 函数有点啰嗦，还需要传一个初始值，
// 在具体的业务函数中，还要关心 result
//total_pay := gReduce(employees, 0.0, func (result float32, e Employee) float32 {
//	return result + e.Salary
//})
//fmt.Printf("Total Salary: %0.2f\n", total_pay) // Total Salary: 43502.05

// refactor: 统计符合条件员工个数
func gCountIf[T any](arr []T, f func(T) bool) int {
	cnt := 0
	for _, elem := range arr {
		if f(elem) {
			cnt += 1
		}
	}
	return cnt
}

// Sumable 接口限定了 U 类型
// 等于多个类型限定
// 语法说明: 只能是 Sumable 里的那些类型，也就是整型或浮点型
//type Sumable interface {
//	type int, int8, int16, int32, int64,
//		uint, uint8, uint16, uint32, uint64,
//		float32, float64
//}

func gSum[T any, U Sumable](arr []T, f func(T) U) U {
	var sum U
	for _, elem := range arr {
		sum += f(elem)
	}
	return sum
}

func main() {
	// 统计年龄大于40岁的员工数
	old := gCountIf(employees, func(e Employee) bool {
		return e.Age > 40
	})
	fmt.Printf("old people(>40): %d\n", old)
	// ld people(>40): 2

	// 统计薪水超过 6000元的员工数
	high_pay := gCountIf(employees, func(e Employee) bool {
		return e.Salary >= 6000
	})
	fmt.Printf("High Salary people(>6k): %d\n", high_pay)
	//High Salary people(>6k): 4

	// 统计年龄小于30岁的员工的薪水
	younger_pay := gSum(employees, func(e Employee) float32 {
		if e.Age < 30 {
			return e.Salary
		}
		return 0
	})
	fmt.Printf("Total Salary of Young People: %0.2f\n", younger_pay)
	//Total Salary of Young People: 19000.00

	// 统计全员的休假天数
	total_vacation := gSum(employees, func(e Employee) int {
		return e.Vacation
	})
	fmt.Printf("Total Vacation: %d day(s)\n", total_vacation)
	//Total Vacation: 32 day(s)

	// 把没有休假的员工过滤出来
	no_vacation := gFilterIn(employees, func(e Employee) bool {
		return e.Vacation == 0
	})
	print(no_vacation)
	//{Hao 44 0 8000.5} {Jack 26 0 4000} {Marry 29 0 6000}
}
