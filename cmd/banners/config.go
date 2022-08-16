package main

import (
	"fmt"
	"time"
)

type Config struct {
	App        App
	Server     Server
	Logger     Logger
	PostgreSQL PostgreSQL
	RabbitMQ   RabbitMQ
}

type App struct {
	Name        string `envconfig:"APP_NAME" default:"banners" required:"true"`
	StorageType string `envconfig:"APP_STORAGE_TYPE" required:"true"`
	RouterType  string `envconfig:"APP_ROUTER_TYPE" required:"true"`
	QueueType   string `envconfig:"APP_QUEUE_TYPE" required:"true"`
}

type Server struct {
	Host         string        `envconfig:"SERVER_HOST" default:"127.0.0.1" required:"true"`
	Port         string        `envconfig:"SERVER_PORT" default:"80" required:"true"`
	ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"30s" required:"true"`
	WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"30s" required:"true"`
	IdleTimeout  time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"30s" required:"true"`
}

type Logger struct {
	Level       string `envconfig:"LOGGER_LEVEL" default:"info" required:"true"`
	Development bool   `envconfig:"LOGGER_DEVELOPMENT" default:"false" required:"true"`
}

type PostgreSQL struct {
	Host     string `envconfig:"PG_HOST" required:"true"`
	Port     int    `envconfig:"PG_PORT" required:"true"`
	Username string `envconfig:"PG_USERNAME" required:"true"`
	Password string `envconfig:"PG_PASSWORD" required:"true"`
	Database string `envconfig:"PG_DATABASE" required:"true"`
}

func (p *PostgreSQL) BuildDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", p.Username, p.Password, p.Host, p.Port, p.Database)
}

type RabbitMQ struct {
	Protocol  string `envconfig:"RABBITMQ_PROTOCOL" required:"true"`
	Username  string `envconfig:"RABBITMQ_USERNAME" required:"true"`
	Password  string `envconfig:"RABBITMQ_PASSWORD" required:"true"`
	Host      string `envconfig:"RABBITMQ_HOST" required:"true"`
	Port      int    `envconfig:"RABBITMQ_PORT" required:"true"`
	Exchange  RabbitMQExchange
	Queue     RabbitMQQueue
	Publisher RabbitMQPublisher
	Consumer  RabbitMQConsumer
}

func (r *RabbitMQ) BuildDSN() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d", r.Protocol, r.Username, r.Password, r.Host, r.Port)
}

type RabbitMQExchange struct {
	Name       string `envconfig:"RABBITMQ_EXCHANGE_NAME" required:"true"`
	Kind       string `envconfig:"RABBITMQ_EXCHANGE_KIND" required:"true"`
	Durable    bool   `envconfig:"RABBITMQ_EXCHANGE_DURABLE" required:"true"`
	AutoDelete bool   `envconfig:"RABBITMQ_EXCHANGE_AUTO_DELETE" required:"true"`
	Internal   bool   `envconfig:"RABBITMQ_EXCHANGE_INTERNAL" required:"true"`
	NoWait     bool   `envconfig:"RABBITMQ_EXCHANGE_NO_WAIT" required:"true"`
}

type RabbitMQQueue struct {
	Name       string `envconfig:"RABBITMQ_QUEUE_NAME" required:"true"`
	Durable    bool   `envconfig:"RABBITMQ_QUEUE_DURABLE" required:"true"`
	AutoDelete bool   `envconfig:"RABBITMQ_QUEUE_AUTO_DELETE" required:"true"`
	Exclusive  bool   `envconfig:"RABBITMQ_QUEUE_EXCLUSIVE" required:"true"`
	NoWait     bool   `envconfig:"RABBITMQ_QUEUE_NO_WAIT" required:"true"`
	BindNoWait bool   `envconfig:"RABBITMQ_QUEUE_BIND_NO_WAIT" required:"true"`
	BindingKey string `envconfig:"RABBITMQ_QUEUE_BINDING_KEY" required:"true"`
}

type RabbitMQPublisher struct {
	Mandatory  bool   `envconfig:"RABBITMQ_PUBLISH_MANDATORY" required:"true"`
	Immediate  bool   `envconfig:"RABBITMQ_PUBLISH_IMMEDIATE" required:"true"`
	RoutingKey string `envconfig:"RABBITMQ_PUBLISH_ROUTING_KEY" required:"true"`
}

type RabbitMQConsumer struct {
	Name      string `envconfig:"RABBITMQ_CONSUMER_NAME" required:"true"`
	AutoAck   bool   `envconfig:"RABBITMQ_CONSUMER_AUTO_ACK" required:"true"`
	Exclusive bool   `envconfig:"RABBITMQ_CONSUMER_EXCLUSIVE" required:"true"`
	NoLocal   bool   `envconfig:"RABBITMQ_CONSUMER_NO_LOCAL" required:"true"`
	NoWait    bool   `envconfig:"RABBITMQ_CONSUMER_NO_WAIT" required:"true"`
}

func NewConfig() Config {
	return Config{}
}
