package server

import (
	"context"
	"fmt"
	"os"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/config"
	"git.edenfarm.id/edenlabs/edenlabs/env"
	"git.edenfarm.id/edenlabs/edenlabs/jwt"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/mongox"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/producer"
	"git.edenfarm.id/edenlabs/edenlabs/provider"
	"git.edenfarm.id/edenlabs/edenlabs/pubsubx"
	"git.edenfarm.id/edenlabs/edenlabs/redisx"
	"git.edenfarm.id/edenlabs/edenlabs/s3x"
	"git.edenfarm.id/edenlabs/edenlabs/telemetry"
	"github.com/evalphobia/logrus_sentry"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func Start() (options opt.Options, err error) {
	env, err := env.Env("env")
	if err != nil {
		fmt.Printf("Failed to using config file, error find %s | %v ", env.ConfigFileUsed(), err)
		return
	}

	cfg, err := config.NewConfig(env)
	if err != nil {
		fmt.Printf("Failed to binding config file, error find %s | %v ", env.ConfigFileUsed(), err)
		return
	}

	app := provider.NewProvider(*cfg)

	// make options for singleton
	options = opt.Options{
		Env:    env,
		Config: *cfg,
	}

	// setup logger
	lgr := logrus.New()
	if cfg.App.Debug {
		lgr.SetFormatter(log.NewFormater(true, cfg.App.Name))
	} else {
		lgr.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}
	lgr.SetReportCaller(true)
	lgr.SetOutput(os.Stdout)

	// setup logger sentry
	if cfg.Sentry.Enabled {
		hook, err := logrus_sentry.NewSentryHook(cfg.Sentry.Dsn, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})
		if err == nil {
			hook.SetEnvironment(cfg.App.Env)
			hook.Timeout = cfg.Sentry.Timeout
			hook.StacktraceConfiguration.Enable = true
			lgr.Hooks.Add(hook)
		}
	}

	logger := log.NewLoggerWithClient(cfg.App.Name, lgr)
	options.Logger = logger

	var dbWrite orm.Ormer
	var dbRead orm.Ormer
	if cfg.Database.Enabled {
		// setup database write
		err = app.GetDBInstanceWrite(logger.Logger())
		if err != nil {
			logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed to start, error connect to database read | %v", err)).Print()
			return
		}
		dbWrite = orm.NewOrmUsingDB("write")

		_, err = orm.GetDB("write")
		if err != nil {
			logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed get database write instance | %v", err)).Print()
			return
		}

		// setup database read
		err = app.GetDBInstanceRead(logger.Logger())
		if err != nil {
			logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed to start, error connect to database write | %v", err)).Print()
			return
		}
		dbRead = orm.NewOrmUsingDB("read")

		_, err = orm.GetDB("read")
		if err != nil {
			logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed get database read instance | %v", err)).Print()
			return
		}
	}
	options.Database = opt.Database{
		Write: dbWrite,
		Read:  dbRead,
	}

	// setup database docs
	var mongoxClient mongox.Client
	if cfg.Mongodb.Enabled {
		var mongoClient *mongo.Client
		mongoClient, err = app.GetDBInstanceDocs()
		if err != nil {
			logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed to start, error connect to database docs | %v", err)).Print()
			return
		}
		mongoxClient = mongox.NewMongox(cfg.Mongodb.Name, mongoClient, logger.Logger())
	}
	options.Mongox = mongoxClient

	// setup database cache
	var redisxClient redisx.Client
	if cfg.Redis.Enabled {
		var redisClient *redis.Client
		redisClient, err = app.GetDBInstanceCache()
		if err != nil {
			logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed to start, error connect to database redis | %v", err)).Print()
			return
		}
		redisxClient = redisx.NewRedisx(redisClient)
	}
	options.Redisx = redisxClient

	// setup database cache
	var s3xClient s3x.Client
	if cfg.S3.Enabled {
		var minioClient *minio.Client
		minioClient, err = app.GetStorageInstanceS3()
		if err != nil {
			logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed to start, error connect to database redis | %v", err)).Print()
			return
		}
		s3xClient = s3x.NewS3x(minioClient)
	}
	options.S3x = s3xClient

	// set up jwt for singleton
	jwtx := jwt.NewJWT([]byte(cfg.Jwt.Key))
	options.Jwt = *jwtx

	// setup trace provider
	var tracer trace.Tracer
	if cfg.Trace.Enabled {
		tp, err := telemetry.NewJaegerTraceProvider(cfg.Trace.Jaeger.Host, cfg.Trace.Jaeger.Port, cfg.App.Name, cfg.App.Env)
		if err != nil {
			logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed get tracer provider instance | %v", err)).Print()
		}
		otel.SetTracerProvider(tp)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		defer func(ctx context.Context) {
			ctx, cancel = context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed to shutdown tracer provider | %v", err)).Print()
			}
		}(ctx)
		tracer = tp.Tracer(cfg.App.Name)
	}
	options.Trace = tracer

	// setup kafka producer
	var kafkaProducer producer.IKafkaProducer
	if cfg.Kafka.Enabled {
		kafkaProducer, err = producer.NewProducer(producer.KafkaProducerOption{
			Brokers: cfg.Kafka.Brokers,
			Topic:   cfg.Kafka.Publisher.Topic,
			Timeout: cfg.Kafka.Timeout,
		})
		if err != nil {
			logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed to connect producer, error connect to kafka producer | %v", err)).Print()
			return
		}
	}

	options.Producer = &producer.Producers{
		KafkaProducer: kafkaProducer,
	}

	// setup pubsub connection
	if cfg.PubSub.Enabled {
		options.PubSub, err = pubsubx.NewPubSub(cfg.PubSub.ID)
		if err != nil {
			logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed to connect pubsub | %v", err)).Print()
		}
	}

	return
}
