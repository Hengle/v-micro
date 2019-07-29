package gotcp

import (
	"fmt"
	"testing"

	"github.com/fananchong/v-micro/transport"
)

func TestGotcpTransport(t *testing.T) {
	tr := NewTransport()

	// bind / listen
	l, err := tr.Listen("localhost:18080")
	if err != nil {
		t.Fatalf("Unexpected error listening %v", err)
	}
	defer l.Close()

	// accept
	go func() {
		if err := l.Accept(func(sock transport.Socket) {
			for {
				var m transport.Message
				if err := sock.Recv(&m); err != nil {
					// 竞态测试时，该协程可能会在测试结束后，执行到这里， t.Logf 会不高兴，panic
					// 因此注释掉 t.Logf ，使用 fmt.Println
					// t.Logf("Exit accept coroutine")
					fmt.Println("Exit accept coroutine")
					return
				}
				t.Logf("Server Received %s", string(m.Body))
				if err := sock.Send(&transport.Message{
					Body: []byte(`pong`),
				}); err != nil {
					return
				}
			}
		}); err != nil {
			t.Fatalf("Unexpected error accepting %v", err)
		}
	}()

	// dial
	c, err := tr.Dial("localhost:18080")
	if err != nil {
		t.Fatalf("Unexpected error dialing %v", err)
	}
	defer c.Close()

	// send <=> receive
	for i := 0; i < 3; i++ {
		if err := c.Send(&transport.Message{
			Body: []byte(`ping`),
		}); err != nil {
			return
		}
		var m transport.Message
		if err := c.Recv(&m); err != nil {
			return
		}
		t.Logf("Client Received %s", string(m.Body))
	}

}

func TestListener(t *testing.T) {
	tr := NewTransport()

	// bind / listen on random port
	l, err := tr.Listen(":0")
	if err != nil {
		t.Fatalf("Unexpected error listening %v", err)
	}
	defer l.Close()

	// try again
	l2, err := tr.Listen(":0")
	if err != nil {
		t.Fatalf("Unexpected error listening %v", err)
	}
	defer l2.Close()

	// now make sure it still fails
	l3, err := tr.Listen(":18080")
	if err != nil {
		t.Fatalf("Unexpected error listening %v", err)
	}
	defer l3.Close()

	if _, err := tr.Listen(":18080"); err == nil {
		t.Fatal("Expected error binding to :18080 got nil")
	}
}
