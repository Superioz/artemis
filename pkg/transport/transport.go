package transport

import (
	"bytes"
	"github.com/satori/go.uuid"
	"github.com/superioz/artemis/pkg/buffer"
	"time"
)

// represents a packet which is being sent to the
// the broker and received from the broker.
type Packet struct {
	// the id of the packet.
	// with this, we can determine, which packet type
	// we received.
	Id uint16

	// the actual data inside the packet.
	// this slice of data is the encoded byte slice
	// of the protobuf messages.
	Data []byte
}

// creates a new packet
func NewPacket(data []byte) (Packet, error) {
	buf := bytes.NewBuffer(data)
	id, err := buffer.ReadUint16(buf)
	if err != nil {
		return Packet{}, err
	}
	return Packet{Id: id, Data: buf.Bytes()}, nil
}

// represents the current state this transporter is in.
// contains all the important information concerning the
// interface and its packet i/o.
type State struct {
	// the unique id of the interface.
	// used to uniquely access interfaces in the
	// broker exchange.
	Id uuid.UUID

	// The key for specifying which exchange to use.
	// This is especially useful for amqp, as `exchange` itself has
	// a specific meaning there. Otherwise, this is just a
	// key to determine which i/o the interface is listening
	// to.
	ExchangeKey string

	// the url of the current broker the `Interface` is
	// connected to.
	// `len(x) = 0` if not connected.
	CurrentBroker string

	// the connection state
	// `true` = is currently connected to a broker.
	Connected bool

	// the timestamp this interface last sent a packet
	// to the broker.
	// `nil` = no outgoing packet yet
	LastSent time.Time

	// the timestamp this interface last received a packet
	// from the broker.
	// `nil` = no incoming packet yet
	LastReceived time.Time
}

// represents the connection interface between the messaging broker
// and this node.
type Interface interface {
	// connects the interface to the broker with `url`
	// also sets the `CurrentBroker` of the `State`
	Connect(url string) error

	// disconnects the interface from the broker
	Disconnect() error

	// the current state of the connection
	// `State#Connected` = false if not connected.
	State() State

	// sends a slice of bytes to given topic.
	// error if data couldn't be sent.
	Send(data []byte, topic string) error

	// returns a channel pipeline for given topic.
	// `nil` if the topic is not registered by this node.
	Receive(topic string) <-chan *Packet
}
