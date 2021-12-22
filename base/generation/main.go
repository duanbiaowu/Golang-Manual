package main

import (
	"fmt"
	"reflect"
)

// 1. 利用 Reflection 实现
type Container struct {
	s reflect.Value
}

func NewContainer(t reflect.Type, size int) *Container {
	if size <= 0 {
		size = 64
	}
	return &Container{
		s: reflect.MakeSlice(reflect.SliceOf(t), 0, size),
	}
}

func (c *Container) Put(val interface{}) error {
	if reflect.ValueOf(val).Type() != c.s.Type().Elem() {
		return fmt.Errorf("Put: cannot put a %T into a slice of %s", val, c.s.Type().Elem())
	}
	c.s = reflect.Append(c.s, reflect.ValueOf(val))
	return nil
}

func (c *Container) Get(val interface{}) error {
	if reflect.ValueOf(val).Kind() != reflect.Ptr || reflect.ValueOf(val).Elem().Type() != c.s.Type().Elem() {
		return fmt.Errorf("Get: needs *%s but got %T", c.s.Type().Elem(), val)
	}
	reflect.ValueOf(val).Elem().Set(c.s.Index(0))
	c.s = c.s.Slice(1, c.s.Len())
	return nil
}

//f1 := 3.1415926
//f2 := 1.1414213
//c := NewContainer(reflect.TypeOf(f1), 16)
//
//if err := c.Put(f1); err != nil {
//	panic(err)
//}
//if err := c.Put(f2); err != nil {
//	panic(err)
//}
//
//g := 0.0
//if err := c.Get(&g); err != nil {
//	panic(err)
//}

// 2. 函数模板，代码生产
// 运行 go generate
//go:generate ./gen.sh ./template/container.tmp.go gen uint32 container
func generateUint32Example() {
	var u uint32 = 42
	c := NewUint32Container()
	c.Put(u)
	v := c.Get()
	fmt.Printf("generateExample: %d (%T)\n", v, v)
}

//go:generate ./gen.sh ./template/container.tmp.go gen string container
func generateStringExample() {
	var s string = "Hello"
	c := NewStringContainer()
	c.Put(s)
	v := c.Get()
	fmt.Printf("generateExample: %s (%T)\n", v, v)
}
