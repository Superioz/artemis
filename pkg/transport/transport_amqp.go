package transport

import (
	"github.com/streadway/amqp"
	"log"
)

type amqpTopic struct {
	key     string
	channel chan *Packet
	queue   amqp.Queue
}

type AMQPInterface struct {
	state *State

	broadcastTopic amqpTopic
	privateTopic   amqpTopic
	leaderTopic    amqpTopic

	connection *amqp.Connection
	channel    *amqp.Channel
}

func (i *AMQPInterface) Connect(url string) error {
	// connects to the amqp broker
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	i.connection = conn

	// TODO
	return nil
}

func (i *AMQPInterface) Disconnect() error {
	return nil
}

func (i *AMQPInterface) State() *State {
	return nil
}

func (i *AMQPInterface) Send(data []byte, topic string) error {
	return nil
}

func (i *AMQPInterface) Receive(topic string) chan *Packet {
	return nil
}
