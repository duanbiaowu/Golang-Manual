//go:build wireinject
// +build wireinject

package struct_providers

import "github.com/google/wire"

//var Set = wire.NewSet(ProvideFoo, ProvideBar, wire.Struct(new(FooBar), "MyFoo", "MyBar"))
var Set = wire.NewSet(ProvideFoo, ProvideBar, wire.Struct(new(FooBar), "*"))

func InitializeFooBar() FooBar {
	wire.Build(Set)
	return FooBar{}
}
