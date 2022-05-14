package creational

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	pre := New()
	cur := New()

	for i := 0; i < 10; i++ {
		assert.Equal(t, pre, cur)
		assert.Same(t, &pre, &cur)
		pre = cur
		cur = New()
	}
}
