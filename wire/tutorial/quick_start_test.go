package tutorial

import "testing"

func TestEvent_DIStart(t *testing.T) {
	message := NewMessage()
	greeter := NewGreeter(message)
	event := NewEvent(greeter)

	event.Start()
}

func TestEvent_WireStart(t *testing.T) {
	event := InitializeEvent()
	event.Start()
}

func TestEvent_WireChangesStart(t *testing.T) {
	event2, err := InitializeEvent2("hello world")
	if err != nil {
		panic(err)
	} else {
		event2.Start()
	}
}
