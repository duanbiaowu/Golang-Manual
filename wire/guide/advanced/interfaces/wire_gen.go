// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package interfaces

import (
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitializeFooBar() string {
	myFoo := providerMyFoo()
	string2 := providerBar(myFoo)
	return string2
}

// wire.go:

var Set = wire.NewSet(providerMyFoo, wire.Bind(new(Foo), new(*MyFoo)), providerBar)
