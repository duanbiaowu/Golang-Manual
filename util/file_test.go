package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteCountBinary(t *testing.T) {
	tests := []struct {
		name string
		args int64
		want string
	}{
		{
			"test-1",
			0,
			"0 B",
		},
		{
			"test-2",
			1023,
			"1023 B",
		},
		{
			"test-3",
			1024,
			"1.0 KB",
		},
		{
			"test-4",
			1024 << 10,
			"1.0 MB",
		},
		{
			"test-5",
			1024 << 20,
			"1.0 GB",
		},
		{
			"test-6",
			1024 << 30,
			"1.0 TB",
		},
		{
			"test-7",
			1024 << 40,
			"1.0 PB",
		},
		{
			"test-8",
			1024 << 50,
			"1.0 EB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ByteCountToReadable(tt.args), "ByteCountToReadable(%v)", tt.args)
		})
	}
}
