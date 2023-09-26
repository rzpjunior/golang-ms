package redisx

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

type Client interface {
	Ping(ctx context.Context) (res string, err error)
	CheckCacheByKey(ctx context.Context, key string) (res bool)
	SetCache(ctx context.Context, key string, value interface{}, expired time.Duration) (err error)
	GetCache(ctx context.Context, key string, value interface{}) (err error)
	DeleteCache(ctx context.Context, key string) (err error)
	DeleteCacheByKey(ctx context.Context, key string) (err error)
	DeleteAll(ctx context.Context) (err error)
	Lock(keyName string) (err error)
	Unlock() (err error)
}

type Redisx struct {
	Client     *redis.Client
	Cache      *cache.Cache
	Redsync    *redsync.Redsync
	Mutex      *redsync.Mutex
	RetryDelay int
	Expiry     int
	RetryCount int
}

func NewRedisx(client *redis.Client) Client {
	redisCache := cache.New(&cache.Options{
		Redis: client,
	})

	return &Redisx{
		Client: client,
		Cache:  redisCache,
	}
}

func (m *Redisx) Ping(ctx context.Context) (res string, err error) {
	res, err = m.Client.Ping(ctx).Result()
	return res, err
}

func (m *Redisx) CheckCacheByKey(ctx context.Context, key string) (res bool) {
	return m.Cache.Exists(ctx, key)
}

func (m *Redisx) SetCache(ctx context.Context, key string, value interface{}, expired time.Duration) (err error) {
	if err = m.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   expired,
	}); err != nil {
		err = fmt.Errorf("failed to set cache | %v", err)
	}
	return
}

func (m *Redisx) GetCache(ctx context.Context, key string, value interface{}) (err error) {
	err = m.Cache.Get(ctx, key, &value)
	if err != nil {
		err = fmt.Errorf("failed to delete cache | %v", err)
	}
	return
}

func (m *Redisx) DeleteCache(ctx context.Context, key string) (err error) {
	err = m.Cache.Delete(ctx, key)
	if err != nil {
		err = fmt.Errorf("failed to delete cache | %v", err)
	}
	return
}

func (m *Redisx) DeleteCacheByKey(ctx context.Context, key string) (err error) {
	if len(os.Args) > 1 {
		key = os.Args[1]
	}
	var foundedRecordCount int = 0
	iter := m.Client.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		m.Cache.Delete(ctx, iter.Val())
		foundedRecordCount++
	}

	if iter.Err() != nil {
		err = fmt.Errorf("failed to delete cache | %v", err)
	}

	return
}

func (m *Redisx) DeleteAll(ctx context.Context) (err error) {
	return m.Client.FlushDB(ctx).Err()
}

func (m *Redisx) Lock(keyName string) (err error) {
	pool := goredis.NewPool(m.Client)
	m.Redsync = redsync.New(pool)
	mutexname := keyName

	if m.RetryDelay == 0 {
		m.RetryDelay = 100
	}
	if m.Expiry == 0 {
		m.Expiry = 60
	}
	if m.RetryCount == 0 {
		m.RetryCount = 64
	}
	m.Mutex = m.Redsync.NewMutex(mutexname, redsync.WithRetryDelay(time.Duration(m.RetryDelay)*time.Millisecond), redsync.WithExpiry(time.Duration(m.Expiry)*time.Second), redsync.WithTries(m.RetryCount))
	if err = m.Mutex.Lock(); err != nil {
		err = errors.New("the process is locked by system")
	}

	return err
}

func (m *Redisx) Unlock() (err error) {
	var statusUnlocked bool
	statusUnlocked, err = m.Mutex.Unlock()
	if !statusUnlocked || err != nil {
		if err != nil {
			err = fmt.Errorf("failed to unlocked the process | %v", err)
			return
		}
	}
	return
}
