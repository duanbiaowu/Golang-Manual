package golang

import "testing"

// go test -run ''      # 执行所有测试。
// go test -run Foo     # 执行匹配 "Foo" 的顶层测试，例如 "TestFooBar"。
// go test -run Foo/A=  # 对于匹配 "Foo" 的顶层测试，执行其匹配 "A=" 的子测试。
// go test -run /A=1    # 执行所有匹配 "A=1" 的子测试。
func TestFoo(t *testing.T) {
	// <setup code>
	t.Run("A=1", func(t *testing.T) {})
	t.Run("A=2", func(t *testing.T) {})
	t.Run("B=1", func(t *testing.T) {})
	// <tear-down code>
}

func TestGroupedParallel(t *testing.T) {
	tests := []struct {
		Name string
	}{
		{
			Name: "test-1",
		},
		{
			Name: "test-2",
		},
	}
	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
		})
	}
}

func TestTeardownParallel(t *testing.T) {
	// This Run will not return until the parallel tests finish.
	t.Run("group", func(t *testing.T) {
		t.Run("Test1", func(t *testing.T) {})
		t.Run("Test2", func(t *testing.T) {})
		t.Run("Test3", func(t *testing.T) {})
	})
	// <tear-down code>
}
