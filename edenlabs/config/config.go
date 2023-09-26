package config

import (
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/env"
	"github.com/go-playground/validator/v10"
)

type (
	// config
	Config struct {
		// app
		App struct {
			Host  string `validate:"required"`
			Port  int    `validate:"required"`
			Name  string `validate:"required"`
			Debug bool   `validate:"required"`
			Env   string `validate:"required"`
		}

		// grpc
		Grpc struct {
			Host string `validate:"required"`
			Port int    `validate:"required"`
		}

		// jwt
		Jwt struct {
			Key string `validate:"required"`
		}

		// database
		Database struct {
			Enabled         bool
			Connection      string        `validate:"required"`
			Name            string        `validate:"required"`
			MaxOpenConns    int           `validate:"required"`
			MaxIdleConns    int           `validate:"required"`
			ConnMaxLifetime time.Duration `validate:"required"`
			// database.write
			Write struct {
				Host     string `validate:"required"`
				Port     int    `validate:"required"`
				Username string `validate:"required"`
				Password string `validate:"required"`
			}
			// database.read
			Read struct {
				Host     string `validate:"required"`
				Port     int    `validate:"required"`
				Username string `validate:"required"`
				Password string `validate:"required"`
			}
		}

		// mongodb
		Mongodb struct {
			Enabled  bool
			Host     string `validate:"required"`
			Port     int    `validate:"required"`
			Name     string `validate:"required"`
			Username string `validate:"required"`
			Password string `validate:"required"`
		}

		// redis
		Redis struct {
			Enabled      bool
			Host         string `validate:"required"`
			Port         int    `validate:"required"`
			Namespace    int
			Username     string
			Password     string
			DialTimeout  time.Duration `validate:"required"`
			WriteTimeout time.Duration `validate:"required"`
			ReadTimeout  time.Duration `validate:"required"`
			IdleTimeout  time.Duration `validate:"required"`
		}

		// sentry
		Sentry struct {
			Enabled bool
			Dsn     string        `validate:"required"`
			Timeout time.Duration `validate:"required"`
		}

		// Trace
		Trace struct {
			Enabled bool
			Jaeger  struct {
				Host string `validate:"required"`
				Port int    `validate:"required"`
			}
		}

		// Metric
		Metric struct {
			Enabled bool
		}

		// client
		Client struct {
			Enabled            bool
			AccountServiceGrpc struct {
				Enabled               bool
				Host                  string        `validate:"required"`
				Port                  int           `validate:"required"`
				Timeout               time.Duration `validate:"required"`
				MaxConcurrentRequests int           `validate:"required"`
				ErrorPercentThreshold int           `validate:"required"`
				Tls                   bool
				PemPath               string
				Secret                string
				Realtime              bool
			}
			AuditServiceGrpc struct {
				Enabled               bool
				Host                  string        `validate:"required"`
				Port                  int           `validate:"required"`
				Timeout               time.Duration `validate:"required"`
				MaxConcurrentRequests int           `validate:"required"`
				ErrorPercentThreshold int           `validate:"required"`
				Tls                   bool
				PemPath               string
				Secret                string
				Realtime              bool
			}
			ConfigurationServiceGrpc struct {
				Enabled               bool
				Host                  string        `validate:"required"`
				Port                  int           `validate:"required"`
				Timeout               time.Duration `validate:"required"`
				MaxConcurrentRequests int           `validate:"required"`
				ErrorPercentThreshold int           `validate:"required"`
				Tls                   bool
				PemPath               string
				Secret                string
				Realtime              bool
			}
			BridgeServiceGrpc struct {
				Enabled               bool
				Host                  string        `validate:"required"`
				Port                  int           `validate:"required"`
				Timeout               time.Duration `validate:"required"`
				MaxConcurrentRequests int           `validate:"required"`
				ErrorPercentThreshold int           `validate:"required"`
				Tls                   bool
				PemPath               string
				Secret                string
				Realtime              bool
			}
			CatalogServiceGrpc struct {
				Enabled               bool
				Host                  string        `validate:"required"`
				Port                  int           `validate:"required"`
				Timeout               time.Duration `validate:"required"`
				MaxConcurrentRequests int           `validate:"required"`
				ErrorPercentThreshold int           `validate:"required"`
				Tls                   bool
				PemPath               string
				Secret                string
				Realtime              bool
			}
			LogisticServiceGrpc struct {
				Enabled               bool
				Host                  string        `validate:"required"`
				Port                  int           `validate:"required"`
				Timeout               time.Duration `validate:"required"`
				MaxConcurrentRequests int           `validate:"required"`
				ErrorPercentThreshold int           `validate:"required"`
				Tls                   bool
				PemPath               string
				Secret                string
				Realtime              bool
			}
			SiteServiceGrpc struct {
				Enabled               bool
				Host                  string        `validate:"required"`
				Port                  int           `validate:"required"`
				Timeout               time.Duration `validate:"required"`
				MaxConcurrentRequests int           `validate:"required"`
				ErrorPercentThreshold int           `validate:"required"`
				Tls                   bool
				PemPath               string
				Secret                string
				Realtime              bool
			}
		}

		// kafka
		Kafka struct {
			Enabled   bool
			Brokers   []string `validate:"required"`
			Version   string
			Username  string
			Password  string
			Tls       bool
			PemPath   string
			Timeout   time.Duration
			Publisher struct {
				Topic string `validate:"required"`
			}
			Subcriber struct {
				Topic string `validate:"required"`
				Group string `validate:"required"`
			}
		}

		// s3
		S3 struct {
			Enabled         bool
			Endpoint        string `validate:"required"`
			BucketName      string `validate:"required"`
			AccessKeyID     string `validate:"required"`
			SecretAccessKey string `validate:"required"`
			Token           string
			UseSSL          bool
		}

		// pubsub
		PubSub struct {
			Enabled bool
			ID      string `validate:"required"`
		}

		// google geocode
		Google struct {
			Enabled    bool
			GeocodeURL string `validate:"required"`
			MapKey     string `validate:"required"`
		}

		// vroom
		Vroom struct {
			Enabled bool
			Url     string `validate:"required"`
		}

		// osrm
		Osrm struct {
			Enabled bool
			Car     string `validate:"required"`
			Bike    string `validate:"required"`
		}

		// print laravel
		PrintService struct {
			Enabled bool
			Url     string `validate:"required"`
		}
	}
)

// NewConfig returns app config.
func NewConfig(env env.Provider) (cfg *Config, err error) {
	cfg = &Config{}
	// validator
	validate := validator.New()

	// app
	cfg.App.Host = env.GetString("app.host")
	cfg.App.Port = env.GetInt("app.port")
	cfg.App.Name = env.GetString("app.name")
	cfg.App.Debug = env.GetBool("app.debug")
	cfg.App.Env = env.GetString("app.env")

	// grpc
	cfg.Grpc.Host = env.GetString("grpc.host")
	cfg.Grpc.Port = env.GetInt("grpc.port")

	// jwt
	cfg.Jwt.Key = env.GetString("jwt.key")

	// database
	cfg.Database.Enabled = env.GetBool("database.enabled")
	cfg.Database.Connection = env.GetString("database.connection")
	cfg.Database.Name = env.GetString("database.name")
	cfg.Database.MaxOpenConns = env.GetInt("database.max_open_conns")
	cfg.Database.MaxIdleConns = env.GetInt("database.max_idle_conns")
	cfg.Database.ConnMaxLifetime = env.GetDuration("database.conn_lifetime_max")
	// database.write
	cfg.Database.Write.Host = env.GetString("database.write.host")
	cfg.Database.Write.Port = env.GetInt("database.write.port")
	cfg.Database.Write.Username = env.GetString("database.write.username")
	cfg.Database.Write.Password = env.GetString("database.write.password")
	// database.read
	cfg.Database.Read.Host = env.GetString("database.read.host")
	cfg.Database.Read.Port = env.GetInt("database.read.port")
	cfg.Database.Read.Username = env.GetString("database.read.username")
	cfg.Database.Read.Password = env.GetString("database.read.password")
	// validate config
	if cfg.Database.Enabled {
		err = validate.Struct(cfg.Database)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// mongodb
	cfg.Mongodb.Enabled = env.GetBool("mongodb.enabled")
	cfg.Mongodb.Host = env.GetString("mongodb.host")
	cfg.Mongodb.Port = env.GetInt("mongodb.port")
	cfg.Mongodb.Name = env.GetString("mongodb.name")
	cfg.Mongodb.Username = env.GetString("mongodb.username")
	cfg.Mongodb.Password = env.GetString("mongodb.password")
	// validate config
	if cfg.Mongodb.Enabled {
		err = validate.Struct(cfg.Mongodb)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// redis
	cfg.Redis.Enabled = env.GetBool("redis.enabled")
	cfg.Redis.Host = env.GetString("redis.host")
	cfg.Redis.Port = env.GetInt("redis.port")
	cfg.Redis.Namespace = env.GetInt("redis.namespace")
	cfg.Redis.Username = env.GetString("redis.username")
	cfg.Redis.Password = env.GetString("redis.password")
	cfg.Redis.DialTimeout = env.GetDuration("redis.dial_timeout")
	cfg.Redis.WriteTimeout = env.GetDuration("redis.write_timeout")
	cfg.Redis.ReadTimeout = env.GetDuration("redis.read_timeout")
	cfg.Redis.IdleTimeout = env.GetDuration("redis.idle_timeout")
	// validate config
	if cfg.Redis.Enabled {
		err = validate.Struct(cfg.Redis)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// sentry
	cfg.Sentry.Enabled = env.GetBool("sentry.enabled")
	cfg.Sentry.Dsn = env.GetString("sentry.dsn")
	cfg.Sentry.Timeout = env.GetDuration("sentry.timeout")
	// validate config
	if cfg.Sentry.Enabled {
		err = validate.Struct(cfg.Sentry)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// trace
	cfg.Trace.Enabled = env.GetBool("trace.enabled")
	// trace.jaeger
	cfg.Trace.Jaeger.Host = env.GetString("trace.jaeger.host")
	cfg.Trace.Jaeger.Port = env.GetInt("trace.jaeger.port")
	// validate config
	if cfg.Trace.Enabled {
		err = validate.Struct(cfg.Trace)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// metric
	cfg.Metric.Enabled = env.GetBool("metric.enabled")

	// client

	// client.account_service_grpc
	cfg.Client.AccountServiceGrpc.Enabled = env.GetBool("client.account_service_grpc.enabled")
	cfg.Client.AccountServiceGrpc.Host = env.GetString("client.account_service_grpc.host")
	cfg.Client.AccountServiceGrpc.Port = env.GetInt("client.account_service_grpc.port")
	cfg.Client.AccountServiceGrpc.Timeout = env.GetDuration("client.account_service_grpc.timeout")
	cfg.Client.AccountServiceGrpc.MaxConcurrentRequests = env.GetInt("client.account_service_grpc.max_concurrent_requests")
	cfg.Client.AccountServiceGrpc.ErrorPercentThreshold = env.GetInt("client.account_service_grpc.error_percent_threshold")
	cfg.Client.AccountServiceGrpc.Tls = env.GetBool("client.account_service_grpc.tls")
	cfg.Client.AccountServiceGrpc.PemPath = env.GetString("client.account_service_grpc.pem_path")
	cfg.Client.AccountServiceGrpc.Secret = env.GetString("client.account_service_grpc.secret")
	cfg.Client.AccountServiceGrpc.Realtime = env.GetBool("client.account_service_grpc.realtime")
	// validate config
	if cfg.Client.AccountServiceGrpc.Enabled {
		err = validate.Struct(cfg.Client.AccountServiceGrpc)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// client.audit_service_grpc
	cfg.Client.AuditServiceGrpc.Enabled = env.GetBool("client.audit_service_grpc.enabled")
	cfg.Client.AuditServiceGrpc.Host = env.GetString("client.audit_service_grpc.host")
	cfg.Client.AuditServiceGrpc.Port = env.GetInt("client.audit_service_grpc.port")
	cfg.Client.AuditServiceGrpc.Timeout = env.GetDuration("client.audit_service_grpc.timeout")
	cfg.Client.AuditServiceGrpc.MaxConcurrentRequests = env.GetInt("client.audit_service_grpc.max_concurrent_requests")
	cfg.Client.AuditServiceGrpc.ErrorPercentThreshold = env.GetInt("client.audit_service_grpc.error_percent_threshold")
	cfg.Client.AuditServiceGrpc.Tls = env.GetBool("client.audit_service_grpc.tls")
	cfg.Client.AuditServiceGrpc.PemPath = env.GetString("client.audit_service_grpc.pem_path")
	cfg.Client.AuditServiceGrpc.Secret = env.GetString("client.audit_service_grpc.secret")
	cfg.Client.AuditServiceGrpc.Realtime = env.GetBool("client.audit_service_grpc.realtime")
	// validate config
	if cfg.Client.AuditServiceGrpc.Enabled {
		err = validate.Struct(cfg.Client.AuditServiceGrpc)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// client.configuration_service_grpc
	cfg.Client.ConfigurationServiceGrpc.Enabled = env.GetBool("client.configuration_service_grpc.enabled")
	cfg.Client.ConfigurationServiceGrpc.Host = env.GetString("client.configuration_service_grpc.host")
	cfg.Client.ConfigurationServiceGrpc.Port = env.GetInt("client.configuration_service_grpc.port")
	cfg.Client.ConfigurationServiceGrpc.Timeout = env.GetDuration("client.configuration_service_grpc.timeout")
	cfg.Client.ConfigurationServiceGrpc.MaxConcurrentRequests = env.GetInt("client.configuration_service_grpc.max_concurrent_requests")
	cfg.Client.ConfigurationServiceGrpc.ErrorPercentThreshold = env.GetInt("client.configuration_service_grpc.error_percent_threshold")
	cfg.Client.ConfigurationServiceGrpc.Tls = env.GetBool("client.configuration_service_grpc.tls")
	cfg.Client.ConfigurationServiceGrpc.PemPath = env.GetString("client.configuration_service_grpc.pem_path")
	cfg.Client.ConfigurationServiceGrpc.Secret = env.GetString("client.configuration_service_grpc.secret")
	cfg.Client.ConfigurationServiceGrpc.Realtime = env.GetBool("client.configuration_service_grpc.realtime")

	// validate config
	if cfg.Client.ConfigurationServiceGrpc.Enabled {
		err = validate.Struct(cfg.Client.ConfigurationServiceGrpc)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// client.bridge_service_grpc
	cfg.Client.BridgeServiceGrpc.Enabled = env.GetBool("client.bridge_service_grpc.enabled")
	cfg.Client.BridgeServiceGrpc.Host = env.GetString("client.bridge_service_grpc.host")
	cfg.Client.BridgeServiceGrpc.Port = env.GetInt("client.bridge_service_grpc.port")
	cfg.Client.BridgeServiceGrpc.Timeout = env.GetDuration("client.bridge_service_grpc.timeout")
	cfg.Client.BridgeServiceGrpc.MaxConcurrentRequests = env.GetInt("client.bridge_service_grpc.max_concurrent_requests")
	cfg.Client.BridgeServiceGrpc.ErrorPercentThreshold = env.GetInt("client.bridge_service_grpc.error_percent_threshold")
	cfg.Client.BridgeServiceGrpc.Tls = env.GetBool("client.bridge_service_grpc.tls")
	cfg.Client.BridgeServiceGrpc.PemPath = env.GetString("client.bridge_service_grpc.pem_path")
	cfg.Client.BridgeServiceGrpc.Secret = env.GetString("client.bridge_service_grpc.secret")
	cfg.Client.BridgeServiceGrpc.Realtime = env.GetBool("client.bridge_service_grpc.realtime")

	// validate config
	if cfg.Client.BridgeServiceGrpc.Enabled {
		err = validate.Struct(cfg.Client.BridgeServiceGrpc)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// client.catalog_service_grpc
	cfg.Client.CatalogServiceGrpc.Enabled = env.GetBool("client.catalog_service_grpc.enabled")
	cfg.Client.CatalogServiceGrpc.Host = env.GetString("client.catalog_service_grpc.host")
	cfg.Client.CatalogServiceGrpc.Port = env.GetInt("client.catalog_service_grpc.port")
	cfg.Client.CatalogServiceGrpc.Timeout = env.GetDuration("client.catalog_service_grpc.timeout")
	cfg.Client.CatalogServiceGrpc.MaxConcurrentRequests = env.GetInt("client.catalog_service_grpc.max_concurrent_requests")
	cfg.Client.CatalogServiceGrpc.ErrorPercentThreshold = env.GetInt("client.catalog_service_grpc.error_percent_threshold")
	cfg.Client.CatalogServiceGrpc.Tls = env.GetBool("client.catalog_service_grpc.tls")
	cfg.Client.CatalogServiceGrpc.PemPath = env.GetString("client.catalog_service_grpc.pem_path")
	cfg.Client.CatalogServiceGrpc.Secret = env.GetString("client.catalog_service_grpc.secret")
	cfg.Client.CatalogServiceGrpc.Realtime = env.GetBool("client.catalog_service_grpc.realtime")

	// client.logistic_service_grpc
	cfg.Client.LogisticServiceGrpc.Enabled = env.GetBool("client.logistic_service_grpc.enabled")
	cfg.Client.LogisticServiceGrpc.Host = env.GetString("client.logistic_service_grpc.host")
	cfg.Client.LogisticServiceGrpc.Port = env.GetInt("client.logistic_service_grpc.port")
	cfg.Client.LogisticServiceGrpc.Timeout = env.GetDuration("client.logistic_service_grpc.timeout")
	cfg.Client.LogisticServiceGrpc.MaxConcurrentRequests = env.GetInt("client.logistic_service_grpc.max_concurrent_requests")
	cfg.Client.LogisticServiceGrpc.ErrorPercentThreshold = env.GetInt("client.logistic_service_grpc.error_percent_threshold")
	cfg.Client.LogisticServiceGrpc.Tls = env.GetBool("client.logistic_service_grpc.tls")
	cfg.Client.LogisticServiceGrpc.PemPath = env.GetString("client.logistic_service_grpc.pem_path")
	cfg.Client.LogisticServiceGrpc.Secret = env.GetString("client.logistic_service_grpc.secret")
	cfg.Client.LogisticServiceGrpc.Realtime = env.GetBool("client.logistic_service_grpc.realtime")

	// client.site_service_grpc
	cfg.Client.SiteServiceGrpc.Enabled = env.GetBool("client.site_service_grpc.enabled")
	cfg.Client.SiteServiceGrpc.Host = env.GetString("client.site_service_grpc.host")
	cfg.Client.SiteServiceGrpc.Port = env.GetInt("client.site_service_grpc.port")
	cfg.Client.SiteServiceGrpc.Timeout = env.GetDuration("client.site_service_grpc.timeout")
	cfg.Client.SiteServiceGrpc.MaxConcurrentRequests = env.GetInt("client.site_service_grpc.max_concurrent_requests")
	cfg.Client.SiteServiceGrpc.ErrorPercentThreshold = env.GetInt("client.site_service_grpc.error_percent_threshold")
	cfg.Client.SiteServiceGrpc.Tls = env.GetBool("client.site_service_grpc.tls")
	cfg.Client.SiteServiceGrpc.PemPath = env.GetString("client.site_service_grpc.pem_path")
	cfg.Client.SiteServiceGrpc.Secret = env.GetString("client.site_service_grpc.secret")
	cfg.Client.SiteServiceGrpc.Realtime = env.GetBool("client.site_service_grpc.realtime")

	// validate config
	if cfg.Client.CatalogServiceGrpc.Enabled {
		err = validate.Struct(cfg.Client.CatalogServiceGrpc)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// kafka
	cfg.Kafka.Enabled = env.GetBool("kafka.enabled")
	cfg.Kafka.Brokers = env.GetStringSlice("kafka.brokers")
	cfg.Kafka.Version = env.GetString("kafka.version")
	cfg.Kafka.Username = env.GetString("kafka.username")
	cfg.Kafka.Password = env.GetString("kafka.password")
	cfg.Kafka.Tls = env.GetBool("kafka.tls")
	cfg.Kafka.PemPath = env.GetString("kafka.pem_tls_path")
	cfg.Kafka.Timeout = env.GetDuration("kafka.timeout")
	// kafka.publisher
	cfg.Kafka.Publisher.Topic = env.GetString("kafka.publisher.topic")
	// kafka.subscriber
	cfg.Kafka.Subcriber.Topic = env.GetString("kafka.subscriber.topic")
	cfg.Kafka.Subcriber.Group = env.GetString("kafka.subscriber.group")
	// validate config
	if cfg.Kafka.Enabled {
		err = validate.Struct(cfg.Kafka)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// s3
	cfg.S3.Enabled = env.GetBool("s3.enabled")
	cfg.S3.Endpoint = env.GetString("s3.endpoint")
	cfg.S3.BucketName = env.GetString("s3.bucket_name")
	cfg.S3.AccessKeyID = env.GetString("s3.access_key_id")
	cfg.S3.SecretAccessKey = env.GetString("s3.secret_access_key")
	cfg.S3.Token = env.GetString("s3.token")
	cfg.S3.UseSSL = env.GetBool("s3.use_ssl")
	// validate config
	if cfg.S3.Enabled {
		err = validate.Struct(cfg.S3)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// pubsub
	cfg.PubSub.Enabled = env.GetBool("pubsub.enabled")
	cfg.PubSub.ID = env.GetString("pubsub.id")
	// validate config
	if cfg.PubSub.Enabled {
		err = validate.Struct(cfg.PubSub)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// google geocode
	cfg.Google.Enabled = env.GetBool("google.enabled")
	cfg.Google.GeocodeURL = env.GetString("google.geocode_url")
	cfg.Google.MapKey = env.GetString("google.map_key")
	// validate config
	if cfg.Google.Enabled {
		err = validate.Struct(cfg.Google)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// vroom
	cfg.Vroom.Enabled = env.GetBool("vroom.enabled")
	cfg.Vroom.Url = env.GetString("vroom.url")
	// validate config
	if cfg.Vroom.Enabled {
		err = validate.Struct(cfg.Vroom)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// osrm
	cfg.Osrm.Enabled = env.GetBool("osrm.enabled")
	cfg.Osrm.Car = env.GetString("osrm.car")
	cfg.Osrm.Bike = env.GetString("osrm.bike")
	// validate config
	if cfg.Osrm.Enabled {
		err = validate.Struct(cfg.Osrm)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	// print laravel
	cfg.PrintService.Enabled = env.GetBool("print_service.enabled")
	cfg.PrintService.Url = env.GetString("print_service.url")
	// validate config
	if cfg.PrintService.Enabled {
		err = validate.Struct(cfg.PrintService)
		if err != nil {
			fmt.Printf("Failed to find specific config file, error find %s | %v ", env.ConfigFileUsed(), err)
			return
		}
	}

	return
}
