package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SimpleVisitor(t *testing.T) {
	c := &Circle{10}
	bytes, err := JsonVisitor(c)
	assert.Nil(t, err)
	assert.Equal(t, `{"Radius":10}`, string(bytes))

	bytes, err = XmlVisitor(c)
	assert.Nil(t, err)
	assert.Equal(t, `<Circle><Radius>10</Radius></Circle>`, string(bytes))
}

// 解耦数据结构和算法
// 使用 Strategy 模式也可以实现，而且更优雅
func Test_SimpleVisitorWithStrategyPattern(t *testing.T) {
	c := Circle{10}
	shapes := []Shape{&c}

	for i := range shapes {
		bytes, err := shapes[i].accept(JsonVisitor)
		assert.Nil(t, err)
		assert.Equal(t, `{"Radius":10}`, string(bytes))

		bytes, err = shapes[i].accept(XmlVisitor)
		assert.Equal(t, `<Circle><Radius>10</Radius></Circle>`, string(bytes))
	}
}

// 以装饰器模式来思考执行流程
func Test_VisitorFunc(t *testing.T) {
	info := Info{}
	var v VisitorCtl = &info
	v = LogVisitor{v}
	v = NameVisitor{v}
	v = OtherThingsVisitor{v}

	loadFile := func(info *Info, err error) error {
		info.Name = "service"
		info.Namespace = "default"
		info.OtherThings = "arguments..."
		return nil
	}
	err := v.Visit(loadFile)
	assert.Nil(t, err)
}

func Test_VisitorFuncWithDecorator(t *testing.T) {
	info := Info{}
	var v VisitorCtl = &info
	v = NewDecoratedVisitorCtl(v, NameVisitorFun, OtherThingsVisitorFun)

	err := v.Visit(func(info *Info, err error) error {
		info.Name = "service"
		info.Namespace = "default"
		info.OtherThings = "arguments..."
		return nil
	})
	assert.Nil(t, err)
}
