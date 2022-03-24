package redis

import (
	"fmt"
	"gin_websocket/lib/config"
	"github.com/go-redis/redis"
	"time"
)

type redisClient struct {
	client *redis.Client
}

var RedisDb redisClient = newClient()

func newClient() redisClient {
	redisConf := config.BaseConf.GetRedisConf()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Password: redisConf.Password,
		DB:       redisConf.Db,
	})
	return redisClient{client: client}
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

func (client *redisClient) HSet(key string, field string, value string) error {
	return client.client.HSet(key, field, value).Err()
}

func (client *redisClient) HScan(key, field string, v interface{}) error {
	return client.client.HGet(key, field).Scan(&v)
}

func (client *redisClient) HVals(key string) ([]string, error) {
	return client.client.HVals(key).Result()
}

func (client *redisClient) HGetAll(key string) (map[string]string, error) {
	return client.client.HGetAll(key).Result()
}

func (client *redisClient) HDelete(key string, field string) error {
	return client.client.HDel(key, field).Err()
}

func (client *redisClient) HExists(key string, field string) (bool, error) {
	return client.client.HExists(key, field).Result()
}
func (client *redisClient) SAdd(key string, value interface{}) error {
	return client.client.SAdd(key, value).Err()
}

//返回个数
func (client *redisClient) SCard(key string) (int, error) {
	count, err := client.client.SCard(key).Result()
	return int(count), err
}

//返回所有成员
func (client *redisClient) SMembers(key string) ([]string, error) {
	return client.client.SMembers(key).Result()
}

//移除成员
func (client *redisClient) SRem(key string, value ...interface{}) error {
	return client.client.SRem(key, value).Err()
}

//是否成员
func (client *redisClient) SIsMember(key string, value interface{}) (bool, error) {
	return client.client.SIsMember(key, value).Result()
}
