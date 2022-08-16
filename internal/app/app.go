package app

import (
	"context"
	"net/http"
	"time"

	"github.com/lozhkindm/otus-banners-rotation/internal/queue"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type Storage interface {
	AddBanner(ctx context.Context, bannerID, slotID int) error
	RemoveBanner(ctx context.Context, bannerID, slotID int) error
	ClickBanner(ctx context.Context, bannerID, slotID, usergroupID int) error
	PickBanner(ctx context.Context, slotID, usergroupID int) (int, error)
}

type Router interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Queue interface {
	DeclareQueue(ctx context.Context, name string) error
	Consume(ctx context.Context, queueName, consumerName string) (<-chan amqp.Delivery, error)
	Close(ctx context.Context) error
	SendEvent(ctx context.Context, queue string, event queue.Event) error
}

type App struct {
	logger      Logger
	storage     Storage
	queue       Queue
	eventsQueue string
}

func New(logger Logger, storage Storage, queue Queue, eventsQueue string) *App {
	return &App{
		logger:      logger,
		storage:     storage,
		queue:       queue,
		eventsQueue: eventsQueue,
	}
}

func (a *App) AddBanner(ctx context.Context, bannerID, slotID int) error {
	return a.storage.AddBanner(ctx, bannerID, slotID)
}

func (a *App) RemoveBanner(ctx context.Context, bannerID, slotID int) error {
	return a.storage.RemoveBanner(ctx, bannerID, slotID)
}

func (a *App) ClickBanner(ctx context.Context, bannerID, slotID, usergroupID int) error {
	if err := a.storage.ClickBanner(ctx, bannerID, slotID, usergroupID); err != nil {
		return err
	}
	event := queue.Event{
		Type:        queue.EventTypeClick,
		BannerID:    bannerID,
		SlotID:      slotID,
		UsergroupID: usergroupID,
		CreatedAt:   time.Now().Unix(),
	}
	if err := a.queue.SendEvent(ctx, a.eventsQueue, event); err != nil {
		a.logger.Error(err.Error())
	}
	return nil
}

func (a *App) PickBanner(ctx context.Context, slotID, usergroupID int) (int, error) {
	bannerID, err := a.storage.PickBanner(ctx, slotID, usergroupID)
	if err != nil {
		return 0, err
	}
	event := queue.Event{
		Type:        queue.EventTypeImpression,
		BannerID:    bannerID,
		SlotID:      slotID,
		UsergroupID: usergroupID,
		CreatedAt:   time.Now().Unix(),
	}
	if err := a.queue.SendEvent(ctx, a.eventsQueue, event); err != nil {
		a.logger.Error(err.Error())
	}
	return bannerID, nil
}
