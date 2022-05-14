package creational

import "sync"

var (
	once     sync.Once
	instance Singleton
)

// Singleton struct
type Singleton struct {
}

// New Singleton
func New() Singleton {
	once.Do(func() {
		instance = Singleton{}
	})
	return instance
}
