package struct_providers

type Foo int
type Bar int

func ProvideFoo() Foo {
	return 0
}

func ProvideBar() Bar {
	return 1
}

type FooBar struct {
	MyFoo Foo
	MyBar Bar
}
