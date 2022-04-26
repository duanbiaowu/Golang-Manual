package testify

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	MockUsers []*User
)

type MockCrawler struct {
	mock.Mock
}

type MockExample struct {
	mock.Mock
}

func init() {
	MockUsers = append(MockUsers, &User{"Java", 30})
	MockUsers = append(MockUsers, &User{"Go", 10})
}

func (m *MockCrawler) GetUserList() ([]*User, error) {
	args := m.Called()
	// a placeholder in the argument list
	return args.Get(0).([]*User), args.Error(1)
}

func (e *MockExample) Hello(n int) int {
	args := e.Mock.Called(n)
	return args.Int(0)
}

func TestAssertSomething(t *testing.T) {
	// assertions equality
	assert.Equal(t, 123, 123, "they should be equal")

	// assertions inequality
	assert.NotEqual(t, 123, 456, "they should not be equal")

	// assertions for nil (good for errors)
	assert.NotNil(t, struct{}{})

	// Every assert func takes the testing.T object as the first argument.
	// This is how it writes the errors out through the normal go test capabilities.
	assertions := assert.New(t)

	assertions.Empty(nil)
	assertions.Empty(0)
	assertions.Empty(0.0)
	assertions.Empty(false)
	assertions.Empty([]byte{})
}

// The require package provides same global functions as the assert package,
// but instead of returning a boolean result they terminate current test.
func TestRequireSomething(t *testing.T) {
	r := require.New(t)

	r.Equal(123, 123)
	//r.Equal(123, 456)
}

func TestMockGetUserList(t *testing.T) {
	// create an instance of our test object
	crawler := new(MockCrawler)

	// setup expectations with a placeholder in the argument list
	crawler.On("GetUserList").Return(MockUsers, nil)

	// call the code we are testing
	GetAndPrintUsers(crawler)

	// assert that the expectations were met
	crawler.AssertExpectations(t)
}

func TestMockExample(t *testing.T) {
	e := new(MockExample)

	e.On("Hello", 1).Return(1).Times(1)
	e.On("Hello", 2).Return(2).Times(2)
	e.On("Hello", 3).Return(3).Times(3)

	ExampleFunc(e)

	e.AssertExpectations(t)
}
