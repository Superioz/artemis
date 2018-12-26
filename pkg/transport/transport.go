package transport

import "time"

type Packet struct {
	Id   uint8
	Data []byte
}

type State struct {
	CurrentBroker string
	Connected     bool
	LastSent      time.Time
	LastReceived  time.Time
}

type Interface interface {
	Connect(url string) error
	Disconnect() error
	State() *State

	Send(data []byte, topic string) error
	Receive() chan *Packet
}
