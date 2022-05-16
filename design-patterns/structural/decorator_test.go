package structural

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_OpDecorator(t *testing.T)  {
	double := func (n int) int {
		return n << 1
	}
	d := OpDecorator(double)
	assert.Equal(t, 10, d(5))

	increment := func(n int) int {
		return n+1
	}
	inc := Operator(increment)
	assert.Equal(t, 6, inc(5))
}