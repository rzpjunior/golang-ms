package db

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type DBRedisOption struct {
	Host         string
	Port         int
	Username     string
	Password     string
	Namespace    int
	DialTimeout  time.Duration
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
}

func NewRedisDatabase(serviceName string, option DBRedisOption) (client *redis.Client, err error) {
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@%s:%d/%d", option.Username, option.Password, option.Host, option.Port, option.Namespace))
	if err != nil {
		err = fmt.Errorf("failed to connect redis %s", err.Error())
		return
	}

	opt.DialTimeout = option.DialTimeout
	opt.WriteTimeout = option.WriteTimeout
	opt.ReadTimeout = option.ReadTimeout
	opt.Password = option.Password
	opt.IdleTimeout = option.IdleTimeout

	client = redis.NewClient(opt)
	return
}
