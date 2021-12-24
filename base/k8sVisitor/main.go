package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// 1. Simple Visitor Example ----------------------------------------
type Visitor func(shape Shape)

type Shape interface {
	accept(Visitor)
}

type Circle struct {
	Radius int
}

func (c *Circle) accept(v Visitor) {
	v(c)
}

type Rectangle struct {
	Width  int
	Height int
}

func (r *Rectangle) accept(v Visitor) {
	v(r)
}

func JsonVisitor(shape Shape) {
	bytes, err := json.Marshal(shape)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

func XmlVisitor(shape Shape) {
	bytes, err := xml.Marshal(shape)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}

// 解耦数据结构和算法
// 使用 Strategy 模式也可以实现，而且更优雅
//c := Circle{10}
//r := Rectangle{100, 50}
//shapes := []Shape{&c, &r}
//
//for i := range shapes {
//	shapes[i].accept(JsonVisitor)
//	shapes[i].accept(XmlVisitor)
//}

// 2. Simplified Kubectl Example ----------------------------------------
type VisitorFunc func(*Info, error) error

type VisitorCtl interface {
	Visit(visitorFunc VisitorFunc) error
}

type Info struct {
	Namespace   string
	Name        string
	OtherThings string
}

func (info *Info) Visit(fn VisitorFunc) error {
	return fn(info, nil)
}

// Name Visitor
type NameVisitor struct {
	visitor VisitorCtl
}

func (v NameVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("NameVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> Name=%s, NameSpace=%s\n", info.Name, info.Namespace)
		}
		fmt.Println("NameVisitor() after call function")
		return err
	})
}

// Other Visitor
type OtherThingsVisitor struct {
	Visitor VisitorCtl
}

func (v OtherThingsVisitor) Visit(fn VisitorFunc) error {
	return v.Visitor.Visit(func(info *Info, err error) error {
		fmt.Println("OtherThingsVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> OtherThings=%s\n", info.OtherThings)
		}
		fmt.Println("OtherThingsVisitor() after call function")
		return err
	})
}

// Log Visitor
type LogVisitor struct {
	visitor VisitorCtl
}

func (v LogVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("LogVisitor() before call function")
		err = fn(info, err)
		fmt.Println("LogVisitor() after call function")
		return err
	})
}

func main() {
	info := Info{}
	var v VisitorCtl = &info
	v = LogVisitor{v}
	v = NameVisitor{v}
	v = OtherThingsVisitor{v}

	loadFile := func(info *Info, err error) error {
		info.Name = "Tom"
		info.Namespace = "NONE"
		info.OtherThings = "We are running as remote team."
		return nil
	}
	_ = v.Visit(loadFile)
}
