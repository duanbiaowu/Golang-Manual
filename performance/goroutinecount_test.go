package performance

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/Jeffail/tunny"
)

func TestControlGoroutineCountUseChannel(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 2)
	for i := 0; i < 10; i++ {
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.Println(i)
			time.Sleep(time.Microsecond)
			<-ch
		}(i)
	}

	wg.Wait()
}

func TestControlGoroutineCountUseTunny(t *testing.T) {
	pool := tunny.NewFunc(2, func(i interface{}) interface{} {
		log.Println(i)
		time.Sleep(time.Microsecond)
		return nil
	})
	defer pool.Close()

	for i := 0; i < 10; i++ {
		go pool.Process(i)
	}
	time.Sleep(time.Second * 1)
}
