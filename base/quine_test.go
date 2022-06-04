package base

import (
	"fmt"
	"testing"
)

func TestQuine(t *testing.T) {
	fmt.Printf("%s%c%s%c\n", q, 0x60, q, 0x60)
}

var q = `
package base

import "fmt"

func TestQuine(t *testing.T) {
    fmt.Printf("%s%c%s%c\n", q, 0x60, q, 0x60)
}

var q = `
