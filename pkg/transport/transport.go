package transport

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
	"github.com/superioz/artemis/pkg/buffer"
	"github.com/superioz/artemis/raft/protocol"
	"time"
)

const (
	// broadcast prefix/suffix
	broadcastKey = "broadcast"

	// the key of the broadcast topic
	BroadcastTopic = broadcastKey + ".*"
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
	Body []byte

	// the whole data set which includes id and content
	Data []byte
}

// creates a new packet
func NewPacket(data []byte) (Packet, error) {
	buf := bytes.NewBuffer(data)
	id, err := buffer.ReadUint16(buf)
	if err != nil {
		return Packet{}, err
	}
	return Packet{Id: id, Body: buf.Bytes(), Data: data}, nil
}

// represents the current state this transporter is in.
// contains all the important information concerning the
// interface and its packet i/o.
type State struct {
	// the unique id of the interface.
	// used to uniquely access interfaces in the
	// broker exchange.
	id uuid.UUID

	// The key for specifying which exchange to use.
	// This is especially useful for amqp, as `exchange` itself has
	// a specific meaning there. Otherwise, this is just a
	// key to determine which i/o the interface is listening
	// to.
	exchangeKey string

	// the url of the current broker the `Interface` is
	// connected to.
	// `len(x) = 0` if not connected.
	currentBroker string

	// the connection state
	// `true` = is currently connected to a broker.
	connected bool

	// the timestamp this interface last sent a packet
	// to the broker.
	// `nil` = no outgoing packet yet
	lastSent time.Time

	// the timestamp this interface last received a packet
	// from the broker.
	// `nil` = no incoming packet yet
	lastReceived time.Time
}

func (s *State) Id() uuid.UUID {
	return s.id
}

func (s *State) CurrentBroker() string {
	return s.currentBroker
}

func (s *State) Connected() bool {
	return s.connected
}

func (s *State) LastSent() time.Time {
	return s.lastSent
}

func (s *State) LastReceived() time.Time {
	return s.lastReceived
}

// represents the connection interface between the messaging broker
// and this node.
type Interface interface {
	// connects the interface to the broker with `url`
	// also handles initialization of the connection important
	// values, which means channel and queue declaring, etc.
	Connect(url string) error

	// disconnects the interface from the broker
	Disconnect() error

	// the current state of the connection
	// `State#Connected` = false if not connected.
	State() *State

	// returns a channel pipeline for this -v
	// sends a slice of bytes to given topic.
	// error if data couldn't be sent.
	Send() chan<- *OutgoingMessage

	// returns a channel pipeline for incoming messages.
	// doesn't matter, from which topic they are.
	Receive() <-chan *IncomingMessage
}

// represents an from amqp received message
type IncomingMessage struct {
	// the raw packet of this per amqp sent message.
	Packet *Packet

	// the topic of the packet.
	Topic string

	// the uuid of the sender, or `nil` if not found.
	Source uuid.UUID

	// the current time stamp when this message
	// got received.
	Time time.Time
}

// represents an outgoing message into the interface
// outgoing channel
type OutgoingMessage struct {
	// the key to route the data to
	RoutingKey string

	// the raw data as byte slice
	Data []byte
}

// unmarshalls given bytes
// fetches the packet id of the data and unmarshals
// the rest of the data into a `proto.Message`
func Decode(data []byte) (proto.Message, error) {
	p, err := NewPacket(data)
	if err != nil {
		return nil, err
	}

	var m = protocol.FromId(p.Id)
	if m == nil {
		return nil, fmt.Errorf("packet with id %d does not exist", p.Id)
	}

	err = proto.Unmarshal(p.Body, m)

	if err != nil {
		return nil, err
	}
	return m, nil
}

// marshals given message into its packet id and message
func Encode(pb proto.Message) ([]byte, error) {
	if pb == nil {
		return nil, fmt.Errorf("message is null")
	}

	packetId := protocol.ToId(pb)
	marshal, err := proto.Marshal(pb)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	buffer.WriteUint16(buf, packetId)
	buf.Write(marshal)

	return buf.Bytes(), nil
}
