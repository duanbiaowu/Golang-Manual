package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"testing"
	"time"
)

func Test_PubSub(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	start := make(chan struct{})
	done := make(chan struct{})
	ctx := context.Background()

	go func() {
		sub := rdb.Subscribe(ctx, "dev")
		ticker := time.NewTicker(time.Second)
		ch := sub.Channel()

		defer func() {
			err := sub.Close()
			if err != nil {
				log.Println(err)
			}
			ticker.Stop()
			done <- struct{}{}
		}()

		start <- struct{}{}

		for {
			select {
			case msg := <-ch:
				fmt.Printf("channel = %s, msg = %s\n", msg.Channel, msg.Payload)
			case <-ticker.C:
				return
			}
		}
	}()

	<-start

	for i := 0; i < 10; i++ {
		err := rdb.Publish(ctx, "dev", i).Err()
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("message [%d] published\n", i)
		}
	}

	<-done
}
