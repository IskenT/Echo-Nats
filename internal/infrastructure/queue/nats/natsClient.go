package natsClient

import (
	"rest_clickhouse/internal/infrastructure/queue"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	Conn *nats.Conn
}

func NewNatsClient(conn *nats.Conn) queue.PubSub {
	return &Nats{Conn: conn}
}

func (ps *Nats) Pub(topic string, data []byte) error {
	return ps.Conn.Publish(topic, data)
}

func (ps *Nats) Sub(topic string, fn func(m *nats.Msg)) (unsub func() error, err error) {
	s, err := ps.Conn.Subscribe(topic, func(msg *nats.Msg) {
		fn(msg)
	})

	if err != nil {
		return nil, err
	}

	return s.Unsubscribe, nil
}
