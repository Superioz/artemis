package transport

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"github.com/superioz/artemis/pkg/logger"
	"strings"
	"time"
)

const (
	// kind of the amqp exchange
	exchangeKind = "topic"

	// the key of the broadcast topic
	broadcastKey = "broadcast"
)

// wrapper struct for an amqp topic.
// contains the channels and information about
// the specific topic.
type amqpTopic struct {
	// the key of the topic.
	// with this key you can route messages to a topic
	// and receive messages from a topic.
	key string

	// the channel for incoming packets
	incoming chan *Packet

	// the consumer channel
	consumer <-chan amqp.Delivery

	// the amqp queue of this topic
	queue amqp.Queue
}

// returns the topic key as lowercase string
// to be able to better compare two keys.
func (t *amqpTopic) LowerKey() string {
	return strings.ToLower(t.key)
}

// implements the transport `Interface` for amqp
type AMQPInterface struct {
	state *State

	// the topic for broadcast messages
	broadcastTopic amqpTopic

	// the topic for private messages explicitly sent to
	// this interface
	privateTopic amqpTopic

	connection  *amqp.Connection
	channel     *amqp.Channel
	notifyClose <-chan *amqp.Error
}

// creates a new amqp interface for connecting
// to e.g. a RabbitMQ broker
func NewAMQPInterface(exchange string) AMQPInterface {
	id := uuid.NewV4()

	return AMQPInterface{
		state:          &State{Id: id, ExchangeKey: exchange},
		broadcastTopic: amqpTopic{key: exchange + "_" + broadcastKey, incoming: make(chan *Packet)},
		privateTopic:   amqpTopic{key: exchange + "_" + id.String(), incoming: make(chan *Packet)},
	}
}

// connects the interface to the amqp broker
func (i *AMQPInterface) Connect(url string) error {
	if i.state.Connected {
		return fmt.Errorf("interface is already connected")
	}

	// connects to the amqp broker
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	i.connection = conn
	i.state.CurrentBroker = url
	i.state.Connected = true
	logger.Info("Connected to amqp.")

	i.notifyClose = conn.NotifyClose(make(chan *amqp.Error))

	// creates unique channel for this interface
	ch, err := conn.Channel()
	if err != nil {
		_ = i.Disconnect()
		return err
	}
	i.channel = ch
	logger.Info("Opened channel to amqp.")

	// check exchange
	err = ch.ExchangeDeclare(i.state.ExchangeKey, exchangeKind, false,
		true, false, false, nil)
	if err != nil {
		_ = i.Disconnect()
		return err
	}
	logger.Info("Declared amqp exchange.")

	// declare the private queue of the interface
	err = i.declareQueue(&i.privateTopic)
	if err != nil {
		_ = i.Disconnect()
		return err
	}

	// declare the broadcast queue
	err = i.declareQueue(&i.broadcastTopic)
	if err != nil {
		_ = i.Disconnect()
		return err
	}
	logger.Info("Declared amqp queues.")

	// listen for input
	go func(i *AMQPInterface) {
		for {
			if !i.state.Connected {
				break
			}

			select {
			case _ = <-i.notifyClose:
				_ = i.Disconnect()
				break
			case broadcast := <-i.broadcastTopic.consumer:
				p, err := NewPacket(broadcast.Body)
				if err != nil {
					logger.Err("couldn't read packet id", err)
					continue
				}

				// update last received timestamp
				i.state.LastReceived = time.Now()

				// send to broadcast consumer if possible
				// TODO not two channels, but one with packets containing the topic
				select {
				case i.broadcastTopic.incoming <- &p:
				}
				break
			case private := <-i.privateTopic.consumer:
				p, err := NewPacket(private.Body)
				if err != nil {
					logger.Err("couldn't read packet id", err)
					continue
				}

				// update last received timestamp
				i.state.LastReceived = time.Now()

				// send to private consumer if possible
				select {
				case i.privateTopic.incoming <- &p:
				}
				break
			}
		}
	}(i)
	return nil
}

// disconnects the interface from the broker
// returns an error, if the connection is already closed
// or if there is any connection issue
func (i *AMQPInterface) Disconnect() error {
	if !i.state.Connected {
		return fmt.Errorf("already disconnected from broker")
	}

	err := i.connection.Close()
	i.state.CurrentBroker = ""
	i.state.Connected = false
	logger.Info("Disconnected from amqp.")
	return err
}

// returns the current interface state
func (i *AMQPInterface) State() State {
	return *i.state
}

// send explicit data to the broker with given `topic`
func (i *AMQPInterface) Send(data []byte, topic string) error {
	if !i.state.Connected {
		return fmt.Errorf("interface is not connected to broker")
	}

	// update last sent
	i.state.LastSent = time.Now()
	err := i.channel.Publish(i.state.ExchangeKey, topic, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	return err
}

// returns the read only packet channel for receiving
// if the topic doesn't exist in this interface return `nil`
func (i *AMQPInterface) Receive(topic string) <-chan *Packet {
	t := i.getTopic(topic)

	if t == nil {
		return nil
	} else {
		return t.incoming
	}
}

// creates and binds a queue for given `topic`.
// returns an error if anything goes wrong, otherwise `nil`.
func (i *AMQPInterface) declareQueue(topic *amqpTopic) error {
	if !i.state.Connected {
		return fmt.Errorf("interface is not connected")
	}
	q, err := i.channel.QueueDeclare(topic.LowerKey(), false, true,
		false, false, nil)
	if err != nil {
		return err
	}
	topic.queue = q

	// binds the queue to the exchange
	err = i.channel.QueueBind(q.Name, topic.LowerKey(), i.state.ExchangeKey, false, nil)
	if err != nil {
		return err
	}

	// get the channel consumer
	topic.consumer, err = i.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

// topic of amqp interface.
// if topic doesn't exist for interface, return `nil`
func (i *AMQPInterface) getTopic(topic string) *amqpTopic {
	var t amqpTopic

	switch strings.ToLower(topic) {
	case i.broadcastTopic.LowerKey():
		t = i.broadcastTopic
		break
	case i.privateTopic.LowerKey():
		t = i.broadcastTopic
		break
	default:
		return nil
	}

	return &t
}
