package xpulsar

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	plog "github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/sirupsen/logrus"

	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/stats"
)

var (
	errNilState        = errors.New("xk6-pubsub: publisher's state is nil")
	errNilStateOfStats = errors.New("xk6-pubsub: stats's state is nil")
)

type PublisherStats struct {
	Topic        string
	ProducerName string
	Messages     int
	Errors       int
	Bytes        int64
}

type PubSub struct{}

func init() {
	modules.Register("k6/x/pulsar", new(PubSub))
}

type PulsarClientConfig struct {
	URL string
}

type ProducerConfig struct {
	Topic string
}

func (p *PubSub) CreateClient(clientConfig PulsarClientConfig) (pulsar.Client, error) {
	logger := logrus.StandardLogger()
	logger.SetLevel(logrus.ErrorLevel)
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               clientConfig.URL,
		ConnectionTimeout: 3 * time.Second,
		Logger:            plog.NewLoggerWithLogrus(logger),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create pulsar client, url: %s, error: %+v", clientConfig.URL, err)
	}
	return client, nil
}

func (p *PubSub) CloseClient(client pulsar.Client) {
	client.Close()
}

func (p *PubSub) CloseProducer(producer pulsar.Producer) {
	producer.Close()
}

func (p *PubSub) CreateProducer(client pulsar.Client, config ProducerConfig) pulsar.Producer {
	option := pulsar.ProducerOptions{
		Topic:               config.Topic,
		Schema:              pulsar.NewStringSchema(nil),
		CompressionType:     pulsar.LZ4,
		CompressionLevel:    pulsar.Faster,
		BatchingMaxMessages: 100,
		MaxPendingMessages:  100,
		SendTimeout:         time.Second,
	}

	producer, err := client.CreateProducer(option)
	if err != nil {
		log.Fatalf("failed to create producer, error: %+v", err)
	}
	return producer
}

func (p *PubSub) Publish(
	ctx context.Context,
	producer pulsar.Producer,
	body []byte,
	properties map[string]string,
	async bool,
) error {
	state := lib.GetState(ctx)
	if state == nil {
		return errNilState
	}

	var err error
	currentStats := PublisherStats{
		Topic:        producer.Topic(),
		ProducerName: producer.Name(),
		Bytes:        int64(len(body)),
		Messages:     1,
	}

	msg := &pulsar.ProducerMessage{
		Value:      "",
		Payload:    body,
		Properties: properties,
	}

	// async send
	if async {
		producer.SendAsync(
			ctx,
			msg,
			func(mi pulsar.MessageID, pm *pulsar.ProducerMessage, e error) {
				if e != nil {
					err = e
					currentStats.Errors++
				}
			},
		)

		return err
	}

	_, err = producer.Send(ctx, msg)
	if err != nil {
		currentStats.Errors++
	}

	if errStats := ReportPubishMetrics(ctx, currentStats); errStats != nil {
		log.Fatal(errStats)
	}

	return err
}

func ReportPubishMetrics(ctx context.Context, currentStats PublisherStats) error {
	state := lib.GetState(ctx)
	if state == nil {
		return errNilStateOfStats
	}

	tags := make(map[string]string)
	tags["producer_name"] = currentStats.ProducerName
	tags["topic"] = currentStats.Topic

	now := time.Now()

	stats.PushIfNotDone(ctx, state.Samples, stats.Sample{
		Time:   now,
		Metric: PublishMessages,
		Tags:   stats.IntoSampleTags(&tags),
		Value:  float64(currentStats.Messages),
	})

	stats.PushIfNotDone(ctx, state.Samples, stats.Sample{
		Time:   now,
		Metric: PublishErrors,
		Tags:   stats.IntoSampleTags(&tags),
		Value:  float64(currentStats.Errors),
	})

	stats.PushIfNotDone(ctx, state.Samples, stats.Sample{
		Time:   now,
		Metric: PublishBytes,
		Tags:   stats.IntoSampleTags(&tags),
		Value:  float64(currentStats.Bytes),
	})
	return nil
}
