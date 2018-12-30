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

// makes sure that packets sent with amqp get
// encoded and decoded properly and consistent.
func TestAMQPPacketConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	exchange := "artemis"
	n1 := NewAMQPInterface(exchange)
	n2 := NewAMQPInterface(exchange)

	group := sync.WaitGroup{}
	mutex := sync.Mutex{}

	received := false
	connected := 0

	run := func(i *AMQPInterface) {
		err := i.Connect("amqp://guest:guest@localhost:5672")

		if err == nil {
			mutex.Lock()
			connected++
			mutex.Unlock()
		}
		group.Done()

		// listens for incoming packet and simply debug a message
		for {
			select {
			case p := <-i.incoming:
				m, err := Decode(p.Packet.Data)
				if err != nil {
					fmt.Println("couldn't decode packet", err)
					break
				}

				switch m.(type) {
				case *protocol.RequestVoteCall:
					// count
					mutex.Lock()
					received = true
					mutex.Unlock()
					break
				default:
					break
				}
				break
			}
		}
	}

	group.Add(2)
	go run(&n1)
	go run(&n2)
	group.Wait()

	if connected < 2 {
		t.Skip("couldn't connect to broker")
	}

	d, _ := Encode(&protocol.RequestVoteCall{Term: 1})

	n1.Send() <- &OutgoingMessage{
		RoutingKey: n2.state.id.String(),
		Data:       d,
	}
	time.Sleep(2 * time.Second)

	if !received {
		t.Error("received packet is not as expected")
	}
	// success
}

// Integration test, if sending a broadcast message works
// with multiple nodes.
// Here we test it with 3 nodes and check if any error
// occurs during the sending of the message
func TestAMQPBroadcast(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	exchange := "artemis"
	n1 := NewAMQPInterface(exchange)

	group := sync.WaitGroup{}
	mutex := sync.Mutex{}

	count := 0
	connected := 0

	run := func(i *AMQPInterface) {
		err := i.Connect("amqp://guest:guest@localhost:5672")

		if err == nil {
			mutex.Lock()
			connected++
			mutex.Unlock()
		}
		group.Done()

		// listens for incoming packet and simply debug a message
		for {
			select {
			case m := <-i.incoming:
				fmt.Println(fmt.Sprintf("{%s, %s, %s, %s}", i.state.id.String(), m.Topic, m.Source, string(m.Packet.Body)))

				// count
				mutex.Lock()
				count++
				mutex.Unlock()
				break
			}
		}
	}

	group.Add(1)
	go run(&n1)

	group.Add(2)
	startNodes(2, exchange, run)
	group.Wait()

	if connected < 3 {
		t.Skip("couldn't connect to broker")
	}

	n1.Send() <- &OutgoingMessage{
		RoutingKey: "broadcast.all",
		Data:       []byte(n1.state.id.String() + ": Hello there!"),
	}
	time.Sleep(2 * time.Second)

	if !n1.State().connected {
		t.Fatal("node disconnected unexpectedly after sending message")
	}
	if count < 3 {
		t.Fatal("didn't receive expected amount of messages")
	}
	// success with sending the message
}

// makes sure that private messages get received by other nodes
func TestAMQPPrivate(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	exchange := "artemis"
	n1 := NewAMQPInterface(exchange)
	n2 := NewAMQPInterface(exchange)

	group := sync.WaitGroup{}
	mutex := sync.Mutex{}

	connected := 0
	count := 0

	run := func(i *AMQPInterface) {
		err := i.Connect("amqp://guest:guest@localhost:5672")

		if err == nil {
			mutex.Lock()
			connected++
			mutex.Unlock()
		}
		group.Done()

		// listens for incoming packet and simply debug a message
		for {
			select {
			case m := <-i.incoming:
				fmt.Println(fmt.Sprintf("{%s, %s, %s, %s}", i.state.id.String(), m.Topic, m.Source, string(m.Packet.Body)))

				// count
				mutex.Lock()
				count++
				mutex.Unlock()
				break
			}
		}
	}

	group.Add(2)
	go run(&n1)
	go run(&n2)
	group.Wait()

	if connected < 2 {
		t.Skip("couldn't connect to broker")
	}
	n1.Send() <- &OutgoingMessage{
		RoutingKey: n2.state.id.String(),
		Data:       []byte(n1.state.id.String() + ": Hello there!"),
	}

	time.Sleep(2 * time.Second)

	if !n1.State().connected {
		t.Fatal("node disconnected unexpectedly after sending message")
	}
	if count != 1 {
		t.Fatal("didn't receive expected amount of messages")
	}
	// success
}

func startNodes(amount int, exchange string, run func(i *AMQPInterface)) {
	for i := 0; i < amount; i++ {
		n := NewAMQPInterface(exchange)

		go run(&n)
	}
}

// makes sure that the registry returns a new instance
// of the stored message and no pointer
func TestRegistry(t *testing.T) {
	m1 := &protocol.RequestVoteCall{Term: 1}
	m2 := &protocol.RequestVoteCall{Term: 2}

	b1, err := Encode(m1)
	b2, err2 := Encode(m2)
	if err != nil || err2 != nil {
		t.Fatal(err)
	}

	m11, err := Decode(b1)
	m21, err2 := Decode(b2)
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

	b, err := Encode(m)
	if err != nil {
		t.Fatal(err)
	}

	m1, err := Decode(b)
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
