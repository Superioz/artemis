package transport

import (
	"bytes"
	"fmt"
	"github.com/superioz/artemis/pkg/buffer"
	"testing"
	"time"
)

func TestNewPacket(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	buffer.WriteUint16(buf, 42)
	buf.WriteString("Hello")

	p, err := NewPacket(buf.Bytes())

	if err != nil {
		t.Error(err)
	}
	if p.Id != 42 {
		t.Error("packet id is not equal to expected id")
	}
}

// Integration test, if sending a broadcast message works
// with multiple nodes.
// Here we test it with 3 nodes and check if any error
// occurs during the sending of the message
func TestAMQPCommunication(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	exchange := "artemis"
	n1 := NewAMQPInterface(exchange)
	n2 := NewAMQPInterface(exchange)
	n3 := NewAMQPInterface(exchange)

	run := func(i *AMQPInterface) {
		err := i.Connect("amqp://guest:guest@localhost:5672")

		if err != nil {
			t.Skip("couldn't connect to test amqp broker")
		}

		for {
			select {
			case p := <-i.privateRoute.consumer:
				fmt.Println(i.state.Id.String()+" | Broadcast | Header:", p.Headers)
				fmt.Println(i.state.Id.String() + " | Private | " + string(p.Body))
				break
			case b := <-i.broadcastRoute.consumer:
				fmt.Println(i.state.Id.String()+" | Broadcast | Header:", b.Headers)
				fmt.Println(i.state.Id.String() + " | Broadcast | " + string(b.Body))
				break
			}
		}
	}

	go run(&n1)
	time.Sleep(2 * time.Second)

	go run(&n2)
	time.Sleep(2 * time.Second)

	go run(&n3)
	time.Sleep(2 * time.Second)

	err := n1.Send([]byte(n1.state.Id.String()+": Hello there!"), "broadcast.all")
	if err != nil {
		t.Error("couldn't send message to exchange")
	}
	// success with sending the message
}
