package transport

import (
	"github.com/satori/go.uuid"
	"time"
)

// represents a packet which is being sent to the
// the broker and received from the broker.
type Packet struct {
	// the id of the packet.
	// with this, we can determine, which packet type
	// we received.
	Id uint8

	// the actual data inside the packet.
	// this slice of data is the encoded byte slice
	// of the protobuf messages.
	Data []byte
}

// represents the current state this transporter is in.
// contains all the important information concerning the
// interface and its packet i/o.
type State struct {
	// the unique id of the interface.
	// used to uniquely access interfaces in the
	// broker exchange.
	Id uuid.UUID

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

type Interface interface {
	Connect(url string) error
	Disconnect() error
	State() *State

	Send(data []byte, topic string) error
	Receive(topic string) chan *Packet
}
