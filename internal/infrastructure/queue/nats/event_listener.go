package natsClient

import (
	"context"
	"encoding/json"
	"rest_clickhouse/internal/infrastructure/queue"
	"rest_clickhouse/internal/infrastructure/usecase/repository"
	"rest_clickhouse/pkg/logger"

	"github.com/nats-io/nats.go"
)

const EventTopicName = "events"

type EventListener struct {
	sub              queue.Subscriber
	eventsRepository repository.EventsRepository
	logger           logger.Logger
	ctx              context.Context
}

func NewEventListener(ctx context.Context, sub queue.Subscriber, logRepository repository.EventsRepository, logger logger.Logger) *EventListener {
	return &EventListener{
		sub:              sub,
		eventsRepository: logRepository,
		logger:           logger,
		ctx:              ctx,
	}
}

func (listen *EventListener) ListenTopic() {
	listen.logger.Info("Event Listener started!")
	unsub, err := listen.sub.Sub(EventTopicName, func(m *nats.Msg) {
		listen.logger.Info("Received a message: %s\n", string(m.Data))

		var itemModel repository.ItemModel
		err := json.Unmarshal(m.Data, &itemModel)

		if err != nil {
			listen.logger.Error(err)
		}

		EventModel := repository.ItemModelToEvent(itemModel)
		err = listen.eventsRepository.Create(EventModel)
		if err != nil {
			listen.logger.Error(err)
		}
	})
	if err != nil {
		listen.logger.ErrorF("Could not subscribe to topic %s: %w", EventTopicName, err)
	}

	go func() {
		<-listen.ctx.Done()
		listen.logger.Info("Stop listen events!")
		if err := unsub(); err != nil {
			panic(err)
		}
	}()
}
