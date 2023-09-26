package server

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"

	"git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/internal/app/handler"
	"git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/internal/app/service"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/sirupsen/logrus"
)

func StartConsumer() (err error) {
	var logger = watermill.NewStdLogger(global.Setup.Common.Config.App.Debug, true)

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		logrus.Errorf("[CONSUMER] Fail to init consumer router : %v", err)
		return
	}

	retryMiddleware := middleware.Retry{
		MaxRetries:      3,
		InitialInterval: time.Millisecond * 10,
	}

	router.AddMiddleware(
		middleware.Recoverer,
		middleware.NewThrottle(10, time.Second).Middleware,
		retryMiddleware.Middleware,
		middleware.CorrelationID,
	)
	router.AddPlugin(plugin.SignalsHandler)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V1_0_0_0
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.ClientID = global.Setup.Common.Config.App.Name

	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               global.Setup.Common.Config.Kafka.Brokers,
			Unmarshaler:           kafka.DefaultMarshaler{},
			ConsumerGroup:         global.Setup.Common.Config.Kafka.Subcriber.Group,
			OverwriteSaramaConfig: saramaConfig,
			OTELEnabled:           false,
			ReconnectRetrySleep:   time.Second * 5,
		},
		logger,
	)
	if err != nil {
		logrus.Errorf("[CONSUMER] Fail to init consumer subscriber : %v", err)
		return
	}

	// setup handler
	boilerplateHandler := handler.BoilerplateConsumerHandler{
		Option:         global.Setup,
		ServicesPerson: service.NewPersonService(),
	}

	router.AddNoPublisherHandler(
		"kafka_consumer",
		global.Setup.Common.Config.Kafka.Subcriber.Topic,
		subscriber,
		boilerplateHandler.IncomingMessage,
	)

	if err = router.Run(context.Background()); err != nil {
		global.Setup.Common.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
