package grpc

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestServerAndClient(t *testing.T) {
	grpcServer := grpc.NewServer()
	RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			t.Error(err)
			return
		}
	}()

	time.Sleep(time.Second)

	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		t.Error(err)
		t.Fatal(err)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			return
		}
	}()

	client := NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &String{Value: "hello"})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello:hello", reply.GetValue())

	// test stream
	// 客户端需要先调用Channel方法获取返回的流对象
	stream, err := client.Channel(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// 在客户端我们将发送和接收操作放到两个独立的Goroutine。
	// 1. 向服务端发送数据：
	go func() {
		for i := 0; i < 5; i++ {
			if err = stream.Send(&String{Value: "hi"}); err != nil {
				t.Error(err)
				return
			}
			time.Sleep(time.Millisecond)
		}
	}()

	// 2.在循环中接收服务端返回的数据：
	for i := 0; i < 5; i++ {
		reply, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		assert.Equal(t, "hello:hi", reply.GetValue())
	}

	time.Sleep(time.Second)
}
