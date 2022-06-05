package grpc

import (
	"context"
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
}
