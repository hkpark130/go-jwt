package adapter

import (
	"github.com/go-redis/redis"
)

func InitializeRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	_, err := client.Ping().Result()

	return client, err
}
