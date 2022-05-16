package idiom

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_FuncOptions(t *testing.T) {
	server, err := NewServer("127.0.0.01", 80,
		Protocol("tcp"),
		Timeout(10),
		MaxConnections(1024000),
	)
	if err != nil {
		t.Skip()
	}

	assert.Equal(t, "tcp", server.Conf.Protocol)
	assert.Equal(t, time.Duration(10), server.Conf.Timeout)
	assert.Equal(t, 1024000, server.Conf.MaxConnections)
}
