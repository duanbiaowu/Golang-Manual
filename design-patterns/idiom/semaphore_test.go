package idiom

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_SemaphoreWithTimeouts(t *testing.T) {
	tickets, timeout := 1, 2*time.Second
	s := NewSem(tickets, timeout)

	err := s.Acquire()
	assert.Nil(t, err)

	err = s.Release()
	assert.Nil(t, err)
}

func Test_SemaphoreWithoutTimeouts(t *testing.T) {
	// 	Non-Blocking
	tickets, timeout := 0, 0*time.Second
	s := NewSem(tickets, timeout)

	err := s.Acquire()
	assert.NotNil(t, err)
	assert.Equal(t, ErrNoTickets, err)
}
