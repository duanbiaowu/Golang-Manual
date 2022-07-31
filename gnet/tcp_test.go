package gnet

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/panjf2000/gnet/v2"
)

type echoServer struct {
	gnet.BuiltinEventEngine

	eng       gnet.Engine
	addr      string
	multicore bool
}

func (es *echoServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	log.Printf("echo server with multi-core=%t is listening on %s\n", es.multicore, es.addr)
	return gnet.None
}

func (es *echoServer) OnTraffic(c gnet.Conn) gnet.Action {
	//buf, _ := c.Next(-1)
	//log.Printf("server read %s", buf)

	buf := make([]byte, 1024)
	read, err := c.Read(buf)
	if err != nil {
		log.Printf("server reading error %s", err)
	}
	log.Printf("server read %d bytes: %s", read, buf)

	_, _ = c.Write(buf)
	return gnet.None
}

func TestGNet(t *testing.T) {
	var port int
	var multicore bool

	// Example command: go run echo.go --port 9000 --multicore=true
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore true")
	flag.Parse()

	echo := &echoServer{addr: fmt.Sprintf("tcp://:%d", port), multicore: multicore}

	// client
	go func() {
		time.Sleep(time.Second)

		dial, err := net.Dial("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			t.Errorf("dial init error %s", err)
			return
		}

		defer func() {
			t.Logf("dial closing...")
			err = dial.Close()
			if err != nil {
				t.Errorf("dial closing error %s", err)
			}
		}()

		go func() {
			data := make([]byte, 1024)
			for {
				read, err := dial.Read(data)
				if err != nil {
					t.Errorf("dial reading error %s", err)
					return
				}
				t.Logf("dial read %d bytes: %s", read, data)
			}
		}()

		for i := 0; i < 10; i++ {
			_, err = dial.Write([]byte("hello world\n"))
			if err != nil {
				t.Errorf("dial writing error %s", err)
			}
			time.Sleep(time.Millisecond)
		}

		time.Sleep(10 * time.Second)
	}()

	log.Fatal(gnet.Run(echo, echo.addr, gnet.WithMulticore(multicore)))
}
