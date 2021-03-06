// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package struct_providers

import (
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitializeFooBar() FooBar {
	foo := ProvideFoo()
	bar := ProvideBar()
	fooBar := FooBar{
		MyFoo: foo,
		MyBar: bar,
	}
	return fooBar
}

func InitializeFooBar2() FooBar2 {
	bar := ProvideBar()
	fooBar2 := FooBar2{
		MyBar: bar,
	}
	return fooBar2
}

// wire.go:

//var Set = wire.NewSet(ProvideFoo, ProvideBar, wire.Struct(new(FooBar), "MyFoo", "MyBar"))
var Set = wire.NewSet(ProvideFoo, ProvideBar, wire.Struct(new(FooBar), "*"))

var Set2 = wire.NewSet(ProvideFoo, ProvideBar, wire.Struct(new(FooBar2), "*"))
