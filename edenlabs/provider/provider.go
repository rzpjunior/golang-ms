package provider

import (
	"git.edenfarm.id/edenlabs/edenlabs/config"
	"git.edenfarm.id/edenlabs/edenlabs/db"
	"git.edenfarm.id/edenlabs/edenlabs/storage"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Provider struct {
	config config.Config
}

// NewProvider initiate Provider object
func NewProvider(config config.Config) *Provider {
	return &Provider{
		config: config,
	}
}

func (a *Provider) GetDBInstanceWrite(lgr *logrus.Logger) (err error) {
	opt := db.DBMysqlOption{
		Host:                 a.config.Database.Write.Host,
		Port:                 a.config.Database.Write.Port,
		Username:             a.config.Database.Write.Username,
		Password:             a.config.Database.Write.Password,
		Name:                 a.config.Database.Name,
		AdditionalParameters: "charset=utf8mb4&collation=utf8mb4_unicode_ci&loc=UTC",
		MaxOpenConns:         a.config.Database.MaxOpenConns,
		MaxIdleConns:         a.config.Database.MaxIdleConns,
		ConnMaxLifetime:      a.config.Database.ConnMaxLifetime,
	}
	err = db.NewMysqlDatabase(a.config.App.Name, lgr, "write", opt)
	return
}

func (a *Provider) GetDBInstanceRead(lgr *logrus.Logger) (err error) {
	opt := db.DBMysqlOption{
		Host:                 a.config.Database.Read.Host,
		Port:                 a.config.Database.Read.Port,
		Username:             a.config.Database.Read.Username,
		Password:             a.config.Database.Read.Password,
		Name:                 a.config.Database.Name,
		AdditionalParameters: "charset=utf8mb4&collation=utf8mb4_unicode_ci&loc=UTC",
		MaxOpenConns:         a.config.Database.MaxOpenConns,
		MaxIdleConns:         a.config.Database.MaxIdleConns,
		ConnMaxLifetime:      a.config.Database.ConnMaxLifetime,
	}
	err = db.NewMysqlDatabase(a.config.App.Name, lgr, "read", opt)
	return
}

func (a *Provider) GetDBInstanceDocs() (client *mongo.Client, err error) {
	opt := db.DBMongoOption{
		Host:     a.config.Mongodb.Host,
		Port:     a.config.Mongodb.Port,
		Username: a.config.Mongodb.Username,
		Password: a.config.Mongodb.Password,
		Name:     a.config.Mongodb.Name,
	}
	client, err = db.NewMongoDatabase(a.config.App.Name, opt)
	return
}

func (a *Provider) GetDBInstanceCache() (client *redis.Client, err error) {
	opt := db.DBRedisOption{
		Host:         a.config.Redis.Host,
		Port:         a.config.Redis.Port,
		Namespace:    a.config.Redis.Namespace,
		Username:     a.config.Redis.Username,
		Password:     a.config.Redis.Password,
		DialTimeout:  a.config.Redis.DialTimeout,
		WriteTimeout: a.config.Redis.WriteTimeout,
		ReadTimeout:  a.config.Redis.ReadTimeout,
		IdleTimeout:  a.config.Redis.IdleTimeout,
	}
	client, err = db.NewRedisDatabase(a.config.App.Name, opt)
	return
}

func (a *Provider) GetStorageInstanceS3() (client *minio.Client, err error) {
	opt := storage.StorageS3Option{
		Endpoint:        a.config.S3.Endpoint,
		AccessKeyID:     a.config.S3.AccessKeyID,
		SecretAccessKey: a.config.S3.SecretAccessKey,
		Token:           a.config.S3.Token,
		UseSSL:          a.config.S3.UseSSL,
	}
	client, err = storage.NewS3Storage(a.config.App.Name, opt)
	return
}
