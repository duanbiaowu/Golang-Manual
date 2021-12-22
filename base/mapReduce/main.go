package main

import (
	"reflect"
)

// Map Example:
func MapStrToStr(arr []string, fn func(s string) string) []string {
	var newArray []string
	for _, val := range arr {
		newArray = append(newArray, fn(val))
	}
	return newArray
}

func MapStrToInt(arr []string, fn func(s string) int) []int {
	var newArray []int
	for _, val := range arr {
		newArray = append(newArray, fn(val))
	}
	return newArray
}

//var names = []string{"Tom", "Terry", "Mary"}
//x := MapStrToStr(names, func(s string) string {
//	return strings.ToLower(s)
//})
//y := MapStrToInt(names, func(s string) int {
//	return len(s)
//})

// Reduce Example:
func Reduce(arr []string, fn func(s string) int) int {
	sum := 0
	for _, val := range arr {
		sum += fn(val)
	}
	return sum
}

//var names = []string{"Tom", "Terry", "Mary"}
//x := Reduce(names, func(s string) int {
//	return len(s)
//})

// Filter Example:
func Filter(arr []int, fn func(n int) bool) []int {
	var newArray []int
	for _, val := range arr {
		if fn(val) {
			newArray = append(newArray, val)
		}
	}
	return newArray
}

//var numbers = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//odds := Filter(numbers, func(n int) bool {
//	return n%2 == 1
//})

// Map/Reduce/Filter只是一种控制逻辑
// 真正的业务逻辑是在传给他们的数据和那个函数来定义的。
// 是的，这是一个很经典的“业务逻辑”和“控制逻辑”分离解耦的编程模式
// 业务逻辑多变，控制逻辑抽象，保持两者分离，然后组合使用，确实经典

// 2. 现实中的业务代码 --------------------------------------------------
type Employee struct {
	Name     string
	Age      int
	Vacation int
	Salary   int
}

var employees = []Employee{
	{"Hao", 44, 0, 8000},
	{"Bob", 34, 10, 5000},
	{"Alice", 23, 5, 9000},
	{"Jack", 26, 0, 4000},
	{"Tom", 48, 9, 7500},
	{"Marry", 29, 0, 6000},
	{"Mike", 32, 8, 4000},
}

func EmployeeCountIf(employees []Employee, fn func(e *Employee) bool) int {
	counts := 0
	// []Employee 使用索引节约内存
	for i, _ := range employees {
		if fn(&employees[i]) {
			counts += 1
		}
	}
	return counts
}

func EmployeeFilterIn(employees []Employee, fn func(e *Employee) bool) []Employee {
	var newList []Employee
	for i, _ := range employees {
		if fn(&employees[i]) {
			newList = append(newList, employees[i])
		}
	}
	return newList
}

func EmployeeSumIf(employees []Employee, fn func(e *Employee) int) int {
	sum := 0
	for i, _ := range employees {
		sum += fn(&employees[i])
	}
	return sum
}

// 统计有多少员工大于40岁
//oldCount := EmployeeCountIf(employees, func(e *Employee) bool {
//	return e.Age > 40
//})
//
//// 统计有多少员工薪水大于6000
//highPayCount := EmployeeCountIf(employees, func(e *Employee) bool {
//	return e.Salary >= 6000
//})
//
//// 统计列出没有休假的员工
//noVacations := EmployeeFilterIn(employees, func(e *Employee) bool {
//	return e.Vacation == 0
//})
//
//// 统计所有员工的薪资总和
//totalPay := EmployeeSumIf(employees, func(e *Employee) int {
//	return e.Salary
//})
//
//// 统计30岁以下员工的薪资总和
//totalYoungPay := EmployeeSumIf(employees, func(e *Employee) int {
//	if e.Age < 30 {
//		return e.Salary
//	}
//	return 0
//})

// 3. 反射基础版本--------------------------------------------------
func Map(data interface{}, fn interface{}) []interface{} {
	vdata := reflect.ValueOf(data)
	vfn := reflect.ValueOf(fn)
	result := make([]interface{}, vdata.Len())

	for i := 0; i < vdata.Len(); i++ {
		result[i] = vfn.Call([]reflect.Value{vdata.Index(i)})[0].Interface()
	}
	return result
}

//square := func(x int) int {
//	return x * x
//}
//nums := []int{1, 2, 3, 4}
//squared := Map(nums, square)
//
//uspace := func(s string) string {
//	return strings.ToUpper(s)
//}
//strs := []string{"Tom", "Terry", "Mary"}
//upStrs := Map(strs, uspace)

// 反射是运行时，如果类型出错，就会报运行时错误
//x := Map(5, 5)
//fmt.Println(x)

// 4. 反射健壮版本--------------------------------------------------
func Transform(slice, fn interface{}) interface{} {
	return transform(slice, fn, false)
}

func TransformInPlace(slice, fn interface{}) interface{} {
	return transform(slice, fn, true)
}

func transform(slice, function interface{}, inPlace bool) interface{} {
	// check slice type is Slice
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("transform: not slice")
	}

	// check the fn signature
	fn := reflect.ValueOf(function)
	elemType := sliceInType.Type().Elem()
	if !verifyFuncSignature(fn, elemType, nil) {
		panic("trasform: function must be of type func(" + sliceInType.Type().Elem().String() + ") outputElemType")
	}

	sliceOutType := sliceInType
	if !inPlace {
		sliceOutType = reflect.MakeSlice(reflect.SliceOf(fn.Type().Out(0)), sliceInType.Len(), sliceInType.Len())
	}
	for i := 0; i < sliceInType.Len(); i++ {
		sliceOutType.Index(i).Set(fn.Call([]reflect.Value{sliceOutType.Index(i)})[0])
	}
	return sliceOutType.Interface()
}

func verifyFuncSignature(fn reflect.Value, types ...reflect.Type) bool {
	// check it is a function
	if fn.Kind() != reflect.Func {
		return false
	}

	// NumIn() - returns a function type's input parameter count.
	// NumOut() - returns a function type's output parameter count.
	if (fn.Type().NumIn() != len(types)-1) || (fn.Type().NumOut() != 1) {
		return false
	}

	// In() - returns the type of a function type's i'th input parameter
	for i := 0; i < len(types)-1; i++ {
		if fn.Type().In(i) != types[i] {
			return false
		}
	}

	// Out() - returns the type of a function type's i'th output parameter.
	outType := types[len(types)-1]
	if outType != nil && fn.Type().Out(0) != outType {
		return false
	}
	return true
}

//// 用于字符串数组
//names := []string{"Tom", "Terry", "Marry"}
//names = Transform(names, func(s string) string {
//	return strings.ToLower(s)
//}).([]string)
//
//// 用于整形数据
//numbers := []int{1, 2, 3, 4, 5}
//TransformInPlace(numbers, func(n int) int {
//	return n * n
//})
//
//// 用于结构体
//employees = []Employee{
//	{"Hao", 44, 0, 8000},
//	{"Bob", 34, 10, 5000},
//	{"Alice", 23, 5, 9000},
//	{"Jack", 26, 0, 4000},
//	{"Tom", 48, 9, 7500},
//}
//TransformInPlace(employees, func(e Employee) Employee {
//	e.Salary += 1000
//	e.Age += 1
//	return e
//})

// 5. Reduce 反射健壮版本--------------------------------------------------
func Reduce2(slice, pairFunc, zero interface{}) interface{} {
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("reduce: wrong type, not slice")
	}

	length := sliceInType.Len()
	if length == 0 {
		return zero
	} else if length == 1 {
		return sliceInType.Index(0)
	}

	elemType := sliceInType.Type().Elem()
	fn := reflect.ValueOf(pairFunc)
	if !verifyFuncSignature(fn, elemType, elemType, elemType) {
		t := elemType.String()
		panic("reduce: function must be of type func(" + t + ", " + t + ") " + t)
	}

	// 除了数组，也可以使用两个变量实现: prev, cur
	var ins [2]reflect.Value
	ins[0] = sliceInType.Index(0)
	ins[1] = sliceInType.Index(1)
	out := fn.Call(ins[:])[0]

	for i := 2; i < length; i++ {
		ins[0] = out
		ins[1] = sliceInType.Index(i)
		out = fn.Call(ins[:])[0]
	}
	return out.Interface()
}

// 6. Filter 反射健壮版本--------------------------------------------------
func Filter2(slice, fn interface{}) interface{} {
	result, _ := filter(slice, fn, false)
	return result
}

func Filter2InPlace(slicePtr, fn interface{}) {
	in := reflect.ValueOf(slicePtr)
	if in.Kind() != reflect.Ptr {
		panic("FilterInPlace: wrong type, " +
			"not a pointer to slice")
	}
	_, n := filter(in.Elem().Interface(), fn, true)
	in.Elem().SetLen(n)
}

var boolType = reflect.ValueOf(true).Type()

func filter(slice, function interface{}, isPlace bool) (interface{}, int) {
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("filter: wrong type, not a slice")
	}

	fn := reflect.ValueOf(function)
	elemType := sliceInType.Type().Elem()
	if !verifyFuncSignature(fn, elemType, boolType) {
		panic("filter: function must be of type func(" + elemType.String() + ") bool")
	}

	var which []int
	for i := 0; i < sliceInType.Len(); i++ {
		if fn.Call([]reflect.Value{sliceInType.Index(i)})[0].Bool() {
			which = append(which, i)
		}
	}

	out := sliceInType
	if !isPlace {
		out = reflect.MakeSlice(sliceInType.Type(), len(which), len(which))
	}
	for i := range which {
		out.Index(i).Set(sliceInType.Index(which[i]))
	}
	return out.Interface(), len(which)
}
