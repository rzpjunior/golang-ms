package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBMongoOption struct {
	Host     string
	Port     int
	Username string
	Password string
	Name     string
}

func NewMongoDatabase(serviceName string, option DBMongoOption) (client *mongo.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d", option.Username, option.Password, option.Host, option.Port)))
	if err != nil {
		err = fmt.Errorf("MongoDB connect %s", err.Error())
		return
	}

	return
}
