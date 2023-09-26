package pubsubx

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type IPubSub interface {
	SendTopic(ctx context.Context, topicID string, msg []byte) (id string, err error)
	ListentoTopic(ctx context.Context, subscriptionID string) *pubsub.Subscription
}

type PubSubOption struct {
	Option *pubsub.Client
}

func NewPubSub(projectID string) (client IPubSub, err error) {
	ctx := context.Background()

	clientPubSub, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return
	}

	client = PubSubOption{
		Option: clientPubSub,
	}

	return
}

func (o PubSubOption) SendTopic(ctx context.Context, topicID string, msg []byte) (id string, err error) {
	var topic *pubsub.Topic

	// create topic or use topic if already exist
	if topic, err = o.Option.CreateTopic(ctx, topicID); err != nil {
		topic = o.Option.Topic(topicID)
		if err != nil {
			return
		}
	}

	result := topic.Publish(ctx, &pubsub.Message{
		Data: msg,
	})

	if id, err = result.Get(ctx); err != nil {
		return
	}

	return
}

func (o PubSubOption) ListentoTopic(ctx context.Context, subscriptionID string) *pubsub.Subscription {
	return o.Option.Subscription(subscriptionID)
}
