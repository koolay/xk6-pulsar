package xpulsar

import "go.k6.io/k6/stats"

var (
	PublishMessages = stats.New("pulsar.publish.message.count", stats.Counter)
	PublishBytes    = stats.New("pulsar.publish.message.bytes", stats.Counter, stats.Data)
	PublishErrors   = stats.New("pulsar.publish.error.count", stats.Counter)
)
