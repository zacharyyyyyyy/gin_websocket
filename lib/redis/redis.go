package redis

import (
	"github.com/go-redis/redis"
	"time"
)

type redisClient struct {
	client *redis.Client
}

var RedisDb redisClient

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "",
		DB:       0,
	})
	RedisDb = redisClient{client: client}
}

func (client *redisClient) Get(key string) (string, error) {
	return client.client.Get(key).Result()

}

func (client *redisClient) Set(key string, value string, time time.Duration) error {
	return client.client.Set(key, value, time).Err()
}

func (client *redisClient) Delete(key string) error {
	return client.client.Del(key).Err()
}

func (client *redisClient) HGet(key string, field string) (string, error) {
	return client.client.HGet(key, field).Result()

}

func (client *redisClient) HSet(key string, field string, value string, time time.Duration) error {
	return client.client.HSet(key, field, value).Err()
}
func (client *redisClient) HScan(key, field string, v interface{}) error {
	return client.client.HGet(key, field).Scan(&v)
}
func (client *redisClient) HDelete(key string, field string) error {
	return client.client.HDel(key, field).Err()
}
func (client *redisClient) HExists(key string, field string) (bool, error) {
	return client.client.HExists(key, field).Result()
}
