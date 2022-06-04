package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTernaryExp(t *testing.T) {
	x, y := 2, 3
	max := If(x > y, x, y).(int)
	assert.Equal(t, 3, max)
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
