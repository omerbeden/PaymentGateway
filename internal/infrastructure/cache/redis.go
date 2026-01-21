package cache

import "github.com/redis/go-redis/v9"

func NewRedis() *redis.Client {
	return &redis.Client{}
}
