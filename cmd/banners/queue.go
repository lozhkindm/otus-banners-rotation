package main

import (
	"context"
	"errors"

	"github.com/lozhkindm/otus-banners-rotation/internal/app"
	rabbitmqqueue "github.com/lozhkindm/otus-banners-rotation/internal/queue/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

const queueTypeRabbitMQ = "rabbitmq"

var undefinedQueueType = errors.New("undefined queue type")

func NewQueue(ctx context.Context, config Config) (app.Queue, error) {
	var queue app.Queue

	switch config.App.QueueType {
	case queueTypeRabbitMQ:
		conn, err := amqp.Dial(config.RabbitMQ.BuildDSN())
		if err != nil {
			return nil, err
		}
		ch, err := conn.Channel()
		if err != nil {
			return nil, err
		}
		rabbit := rabbitmqqueue.New(conn, ch)
		rabbit.QueueSettings = rabbitmqqueue.QueueSettings{
			Durable:    config.RabbitMQ.Queue.Durable,
			AutoDelete: config.RabbitMQ.Queue.AutoDelete,
			Exclusive:  config.RabbitMQ.Queue.Exclusive,
			NoWait:     config.RabbitMQ.Queue.NoWait,
			BindNoWait: config.RabbitMQ.Queue.BindNoWait,
			BindingKey: config.RabbitMQ.Queue.BindingKey,
		}
		rabbit.PublisherSettings = rabbitmqqueue.PublisherSettings{
			Mandatory:  config.RabbitMQ.Publisher.Mandatory,
			Immediate:  config.RabbitMQ.Publisher.Immediate,
			RoutingKey: config.RabbitMQ.Publisher.RoutingKey,
		}
		rabbit.ExchangeSettings = rabbitmqqueue.ExchangeSettings{
			Kind:       config.RabbitMQ.Exchange.Kind,
			Durable:    config.RabbitMQ.Exchange.Durable,
			AutoDelete: config.RabbitMQ.Exchange.AutoDelete,
			Internal:   config.RabbitMQ.Exchange.Internal,
			NoWait:     config.RabbitMQ.Exchange.NoWait,
		}
		rabbit.ConsumerSettings = rabbitmqqueue.ConsumerSettings{
			AutoAck:   config.RabbitMQ.Consumer.AutoAck,
			Exclusive: config.RabbitMQ.Consumer.Exclusive,
			NoLocal:   config.RabbitMQ.Consumer.NoLocal,
			NoWait:    config.RabbitMQ.Consumer.NoWait,
		}
		if err := rabbit.DeclareExchange(ctx, config.RabbitMQ.Exchange.Name); err != nil {
			return nil, err
		}
		queue = rabbit
	default:
		return nil, undefinedQueueType
	}

	return queue, nil
}
