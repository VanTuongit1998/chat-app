package configs

import (
	"github.com/go-redis/redis/v8"
)

func NewClient(address string) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: address})
}
