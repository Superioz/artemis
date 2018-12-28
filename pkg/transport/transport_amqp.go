package transport

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"github.com/superioz/artemis/pkg/logger"
	"time"
)

const (
	// kind of the amqp exchange
	exchangeKind = "topic"
)

// wrapper struct for an amqp topic.
// contains the channels and information about
// the specific topic.
type amqpRoute struct {
	// the key of the topic.
	// with this key you can route messages to a topic
	// and receive messages from a topic.
	topic string

	// name of the queue to be created
	queueName string

	// the consumer channel
	consumer <-chan amqp.Delivery

	// the amqp queue of this topic
	queue amqp.Queue
}

// implements the transport `Interface` for amqp
type AMQPInterface struct {
	state *State

	// the route for broadcast messages
	broadcastRoute amqpRoute

	// the route for private messages explicitly sent to
	// this interface
	privateRoute amqpRoute

	// the messaging channel for incoming messages
	incoming chan *IncomingMessage

	// the messaging channel for outgoing messages
	outgoing chan *OutgoingMessage

	connection  *amqp.Connection
	channel     *amqp.Channel
	notifyClose <-chan *amqp.Error
}

// creates a new amqp interface for connecting
// to e.g. a RabbitMQ broker
func NewAMQPInterface(exchange string) AMQPInterface {
	id := uuid.NewV4()

	return AMQPInterface{
		state: &State{Id: id, ExchangeKey: exchange},
		broadcastRoute: amqpRoute{
			topic:     BroadcastTopic,
			queueName: fmt.Sprintf("%s_%s_%s", exchange, id.String(), broadcastKey),
		},
		privateRoute: amqpRoute{
			topic:     id.String(),
			queueName: fmt.Sprintf("%s_%s", exchange, id.String()),
		},
	}
}

// connects the interface to the amqp broker
// if the connection somehow fails, `Disconnect` will be
// executed. It is neccessary to execute this method again to
// create the needed channels and to initialize important
// other parts of the connection such as the queues and the
// exchange from amqp.
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
	err = i.declareQueue(&i.privateRoute)
	if err != nil {
		_ = i.Disconnect()
		return err
	}

	// declare the broadcast queue
	err = i.declareQueue(&i.broadcastRoute)
	if err != nil {
		_ = i.Disconnect()
		return err
	}
	logger.Info("Declared amqp queues.")

	// create channels
	i.incoming = make(chan *IncomingMessage)
	i.outgoing = make(chan *OutgoingMessage)

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
			case broadcast := <-i.broadcastRoute.consumer:
				m, err := convertMessage(broadcast, i.broadcastRoute)
				if err != nil {
					logger.Err("couldn't read message", err)
					continue
				}

				// update last received timestamp
				i.state.LastReceived = m.Time

				// send to incoming if possible
				select {
				case i.incoming <- &m:
				}
				break
			case private := <-i.privateRoute.consumer:
				m, err := convertMessage(private, i.privateRoute)
				if err != nil {
					logger.Err("couldn't read message", err)
					continue
				}

				// update last received timestamp
				i.state.LastReceived = m.Time

				// send to incoming if possible
				select {
				case i.incoming <- &m:
				}
				break
			}
		}
	}(i)

	// listen for output
	go func(i *AMQPInterface) {
		for outgoing := range i.outgoing {
			if !i.state.Connected {
				break
			}

			// header
			table := amqp.Table{}
			table["user-id"] = i.state.Id.String()

			// update last sent
			i.state.LastSent = time.Now()
			err := i.channel.Publish(i.state.ExchangeKey, outgoing.RoutingKey, false, false,
				amqp.Publishing{
					Headers:     table,
					ContentType: "text/plain",
					Body:        outgoing.Data,
				})
			if err != nil {
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
	close(i.incoming)
	close(i.outgoing)
	logger.Info("Disconnected from amqp.")
	return err
}

// returns the current interface state
func (i *AMQPInterface) State() State {
	return *i.state
}

// returns the write only packet channel for sending.
func (i *AMQPInterface) Send() chan<- *OutgoingMessage {
	return i.outgoing
}

// returns the read only packet channel for receiving
func (i *AMQPInterface) Receive() <-chan *IncomingMessage {
	return i.incoming
}

// converts received `Delivery` into a wrapper message struct.
// stores the current timestamp, the route and the source, which
// a normal `Packet` doesn't.
func convertMessage(d amqp.Delivery, route amqpRoute) (IncomingMessage, error) {
	p, err := NewPacket(d.Body)
	if err != nil {
		return IncomingMessage{}, err
	}

	// create wrapper for this message
	uidStr, ok := d.Headers["user-id"].(string)
	var uid uuid.UUID
	if ok {
		uid, _ = uuid.FromString(uidStr)
	}

	message := IncomingMessage{
		Packet: &p,
		Topic:  route.topic,
		Source: uid,
		Time:   time.Now(),
	}
	return message, nil
}

// creates and binds a queue for given `topic`.
// returns an error if anything goes wrong, otherwise `nil`.
func (i *AMQPInterface) declareQueue(topic *amqpRoute) error {
	if !i.state.Connected {
		return fmt.Errorf("interface is not connected")
	}
	q, err := i.channel.QueueDeclare(topic.queueName, false, true,
		false, false, nil)
	if err != nil {
		return err
	}
	topic.queue = q

	// binds the queue to the exchange
	err = i.channel.QueueBind(q.Name, topic.topic, i.state.ExchangeKey, false, nil)
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
