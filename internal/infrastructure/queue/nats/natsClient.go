package natsClient

import (
	"rest_clickhouse/internal/infrastructure/queue"

	"github.com/nats-io/nats.go"
)

// Nats реализует интерфейс PubSub для взаимодействия с NATS.
type Nats struct {
	Conn *nats.Conn
}

// NewNatsClient создает новый клиент NATS на основе переданного соединения.
func NewNatsClient(conn *nats.Conn) queue.PubSub {
	return &Nats{Conn: conn}
}

// Pub публикует сообщение в указанный топик NATS.
func (n *Nats) Pub(topic string, data []byte) error {
	return n.Conn.Publish(topic, data)
}

// Sub подписывается на сообщения в указанном топике NATS и вызывает функцию обратного вызова для каждого полученного сообщения.
// Возвращает функцию для отписки от топика и ошибку, если подписка не удалась.
func (n *Nats) Sub(topic string, fn func(m *nats.Msg)) (unsub func() error, err error) {
	sub, err := n.Conn.Subscribe(topic, func(msg *nats.Msg) {
		fn(msg)
	})
	if err != nil {
		return nil, err
	}

	return sub.Unsubscribe, nil
}
