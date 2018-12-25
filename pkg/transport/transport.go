package transport

type Packet struct {
	Id uint8
	Data []byte
}

type State struct {
	Connected bool
}

type Interface interface {
	Connect(url string)
	Disconnect()
	State() State

	Send(data []byte, topic string) error
	Receive() chan Packet
}
