package clientx

import (
	"time"

	"github.com/afex/hystrix-go/hystrix"
	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"
	"github.com/gojek/heimdall/v7"
)

// NewGRPCClient returns a new instance of hystrix grpc Client
func NewGRPCClient(opts ...Option) *Client {
	client := Client{
		timeout:                defaultHTTPTimeout,
		hystrixTimeout:         defaultHystrixTimeout,
		maxConcurrentRequests:  defaultMaxConcurrentRequests,
		errorPercentThreshold:  defaultErrorPercentThreshold,
		sleepWindow:            defaultSleepWindow,
		requestVolumeThreshold: defaultRequestVolumeThreshold,
		retryCount:             defaultHystrixRetryCount,
		retrier:                heimdall.NewNoRetrier(),
	}

	for _, opt := range opts {
		opt(&client)
	}

	if client.statsD != nil {
		c, err := plugins.InitializeStatsdCollector(client.statsD)
		if err != nil {
			panic(err)
		}

		metricCollector.Registry.Register(c.NewStatsdCollector)
	}

	hystrix.ConfigureCommand(client.hystrixCommandName, hystrix.CommandConfig{
		Timeout:                durationToInt(client.hystrixTimeout, time.Millisecond),
		MaxConcurrentRequests:  client.maxConcurrentRequests,
		RequestVolumeThreshold: client.requestVolumeThreshold,
		SleepWindow:            client.sleepWindow,
		ErrorPercentThreshold:  client.errorPercentThreshold,
	})

	return &client
}

// Execute makes an GRPC request
func (hhc *Client) Execute(f func() error) (err error) {
	for i := 0; i <= hhc.retryCount; i++ {
		err = hystrix.Do(hhc.hystrixCommandName, f, hhc.fallbackFunc)

		if err != nil {
			backoffTime := hhc.retrier.NextInterval(i)
			time.Sleep(backoffTime)
			continue
		}

		break
	}

	return err
}
