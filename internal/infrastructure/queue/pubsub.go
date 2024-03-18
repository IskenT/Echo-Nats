package queue

import "github.com/nats-io/nats.go"

// Publisher определяет интерфейс для публикации сообщений.
type Publisher interface {
	Pub(topic string, data []byte) error
}

// Subscriber определяет интерфейс для подписки на сообщения.
type Subscriber interface {
	Sub(topic string, fn func(m *nats.Msg)) (unsub func() error, err error)
}

// PubSub объединяет интерфейсы Publisher и Subscriber.
type PubSub interface {
	Publisher
	Subscriber
}
