package queue

import "github.com/nats-io/nats.go"

type Publisher interface {
	Pub(topic string, data []byte) error
}

type Subscriber interface {
	Sub(topic string, fn func(m *nats.Msg)) (unsub func() error, err error)
}

type PubSub interface {
	Publisher
	Subscriber
}
