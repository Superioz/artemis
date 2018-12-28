package transport

import (
	"bytes"
	"fmt"
	"github.com/superioz/artemis/pkg/buffer"
	"github.com/superioz/artemis/raft/protocol"
	"sync"
	"testing"
	"time"
)

// tests if a created packet contains expected values
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

	mutex := sync.Mutex{}
	count := 0

	run := func(i *AMQPInterface) {
		err := i.Connect("amqp://guest:guest@localhost:5672")

		if err != nil {
			t.Skip("couldn't connect to test amqp broker")
		}

		// listens for incoming packet and simply debug a message
		for {
			select {
			case m := <-i.incoming:
				fmt.Println(fmt.Sprintf("{%s, %s, %s, %s}", i.state.Id.String(), m.Topic, m.Source, string(m.Packet.Data)))

				// count
				mutex.Lock()
				count++
				mutex.Unlock()
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

	n1.Send() <- &OutgoingMessage{
		RoutingKey: "broadcast.all",
		Data:       []byte(n1.state.Id.String() + ": Hello there!"),
	}
	time.Sleep(2 * time.Second)

	if !n1.State().Connected {
		t.Fatal("node disconnected unexpectedly after sending message")
	}
	if count < 3 {
		t.Fatal("didn't receive expected amount of messages")
	}
	// success with sending the message
}

// makes sure that the registry returns a new instance
// of the stored message and no pointer
func TestRegistry(t *testing.T) {
	m1 := &protocol.RequestVoteCall{Term: 1}
	m2 := &protocol.RequestVoteCall{Term: 2}

	b1, err := Marshal(m1)
	b2, err2 := Marshal(m2)
	if err != nil || err2 != nil {
		t.Fatal(err)
	}

	m11, err := Unmarshal(b1)
	m21, err2 := Unmarshal(b2)
	if err != nil || err2 != nil {
		t.Fatal(err)
	}

	origin1, _ := m11.(*protocol.RequestVoteCall)
	term1 := origin1.Term

	origin2, _ := m21.(*protocol.RequestVoteCall)
	term2 := origin2.Term

	if term1 == term2 {
		t.Fatal("found identical but expected different packets")
	}
}

// makes sure that the marshalling works correctly
func TestMarshal(t *testing.T) {
	m := &protocol.RequestVoteCall{Term: 1}

	b, err := Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	m1, err := Unmarshal(b)
	if err != nil {
		t.Fatal(err)
	}

	origin, ok := m1.(*protocol.RequestVoteCall)
	if !ok {
		t.Fatal("couldn't convert to original type")
	}
	if origin.Term != m.Term {
		t.Fatal("expected term is not real term")
	}
}
