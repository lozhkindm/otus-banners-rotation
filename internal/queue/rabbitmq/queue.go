package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/lozhkindm/otus-banners-rotation/internal/queue"

	amqp "github.com/rabbitmq/amqp091-go"
)

var errTimeoutClosingRabbitMQ = errors.New("timeout closing rabbitmq")

type ExchangeSettings struct {
	Kind                                  string
	Durable, AutoDelete, Internal, NoWait bool
}

type QueueSettings struct {
	Durable, AutoDelete, Exclusive, NoWait, BindNoWait bool
	BindingKey                                         string
}

type PublisherSettings struct {
	Mandatory, Immediate bool
	RoutingKey           string
}

type ConsumerSettings struct {
	AutoAck, Exclusive, NoLocal, NoWait bool
}

type Queue struct {
	Queue             amqp.Queue
	Conn              *amqp.Connection
	Ch                *amqp.Channel
	ExchangeSettings  ExchangeSettings
	QueueSettings     QueueSettings
	PublisherSettings PublisherSettings
	ConsumerSettings  ConsumerSettings
}

func New(conn *amqp.Connection, ch *amqp.Channel) Queue {
	return Queue{
		Conn: conn,
		Ch:   ch,
	}
}

func (q Queue) DeclareExchange(_ context.Context, name string) error {
	return q.Ch.ExchangeDeclare(
		name,
		q.ExchangeSettings.Kind,
		q.ExchangeSettings.Durable,
		q.ExchangeSettings.AutoDelete,
		q.ExchangeSettings.Internal,
		q.ExchangeSettings.NoWait,
		nil,
	)
}

func (q Queue) DeclareQueue(_ context.Context, name string) (err error) {
	q.Queue, err = q.Ch.QueueDeclare(
		name,
		q.QueueSettings.Durable,
		q.QueueSettings.AutoDelete,
		q.QueueSettings.Exclusive,
		q.QueueSettings.NoWait,
		nil,
	)
	return
}

func (q Queue) BindQueue(_ context.Context, queueName, exchangeName string) error {
	return q.Ch.QueueBind(
		queueName,
		q.QueueSettings.BindingKey,
		exchangeName,
		q.QueueSettings.BindNoWait,
		nil,
	)
}

func (q Queue) Consume(_ context.Context, queueName, consumerName string) (<-chan amqp.Delivery, error) {
	return q.Ch.Consume(
		queueName,
		consumerName,
		q.ConsumerSettings.AutoAck,
		q.ConsumerSettings.Exclusive,
		q.ConsumerSettings.NoLocal,
		q.ConsumerSettings.NoWait,
		nil,
	)
}

func (q Queue) SendEvent(ctx context.Context, queue string, event queue.Event) (err error) {
	var bytes []byte
	if bytes, err = json.Marshal(&event); err != nil {
		return err
	}
	return q.Ch.PublishWithContext(
		ctx,
		queue,
		q.PublisherSettings.RoutingKey,
		q.PublisherSettings.Mandatory,
		q.PublisherSettings.Immediate,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		})
}

func (q Queue) Close(ctx context.Context) error {
	var (
		cnt            = 0
		ch             = make(chan error, 2)
		newCtx, cancel = context.WithTimeout(ctx, 3*time.Second)
	)
	defer cancel()

	go func() {
		ch <- q.Ch.Close()
		ch <- q.Conn.Close()
	}()

	for {
		select {
		case <-newCtx.Done():
			return errTimeoutClosingRabbitMQ
		case err := <-ch:
			if err != nil {
				return err
			}
			cnt++
			if cnt == 2 {
				return err
			}
		}
	}
}
