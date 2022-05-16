package idiom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Generate(t *testing.T) {
	gen := Generate(1, 100)
	cnt := 1
	for i := range gen {
		assert.Equal(t, cnt, i)
		cnt++
		if i == 100 {
			break
		}
	}
	close(gen)
}
