package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type IKafkaProducer interface {
	PublishMessage(ctx context.Context, messagePayload interface{}) (err error)
}

type KafkaProducerOption struct {
	Brokers []string
	Topic   string
	Timeout time.Duration
}

type kafkaProducer struct {
	Option    KafkaProducerOption
	Publisher *kafka.Publisher
}

func NewProducer(opt KafkaProducerOption) (iKafkaProducer IKafkaProducer, err error) {
	logger := watermill.NewStdLogger(false, false)

	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   opt.Brokers,
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		return
	}

	iKafkaProducer = kafkaProducer{
		Option:    opt,
		Publisher: kafkaPublisher,
	}

	return
}

func (o kafkaProducer) PublishMessage(ctx context.Context, messagePayload interface{}) (err error) {
	payload, err := json.Marshal(messagePayload)
	if err != nil {
		return
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)

	middleware.SetCorrelationID(watermill.NewShortUUID(), msg)
	middleware.Timeout(o.Option.Timeout)
	middleware.NewThrottle(5, time.Second)

	topic := o.Option.Topic

	err = o.Publisher.Publish(topic, msg)
	if err != nil {
		err = fmt.Errorf("failed publish to kafka | %v", err)
		return
	}

	return
}
