package redis

import (
	"context"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	sleep  time.Duration
}

func NewClient(sleep time.Duration) *RedisClient {
	client, _ := redismock.NewClientMock()

	return &RedisClient{
		client: client,
		sleep:  sleep,
	}
}

func (rc *RedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	// mock latency
	time.Sleep(rc.sleep)
	return rc.client.Incr(ctx, key)
}

func (rc *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return rc.client.Get(ctx, key)
}
