package redis

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewClient() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")

	client := redis.NewClient(&redis.Options{Addr: addr})

	return client
}
