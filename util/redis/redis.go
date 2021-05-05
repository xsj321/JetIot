package redis

import (
	"JetIot/conf"
	"github.com/go-redis/redis"
	"time"
)

var redisDb *redis.Client

const (
	FRIEND_ALL  = "all_user_list"
	FRIEND_LIST = "friend_list"
	GROUP_LIST  = "group_list"
)

func InitRedis() {
	redisDb = redis.NewClient(&redis.Options{
		Addr:     conf.Default.RedisServer + ":" + conf.Default.RedisPort, // use default Addr
		Password: conf.Default.RedisPassword,                              // no password set
		DB:       0,                                                       // use default DB
	})
}

func Set(key string, value interface{}) error {
	err := redisDb.Set(key, value, 0).Err()
	return err
}

func SetWithExpiration(key string, value interface{}, exp time.Duration) error {
	err := redisDb.Set(key, value, time.Second*exp).Err()
	return err
}

func Get(key string) ([]byte, error) {
	get := redisDb.Get(key)
	return get.Bytes()
}

func Del(key string) error {
	err := redisDb.Del(key).Err()
	return err
}

func SAdd(key string, value interface{}) {
	redisDb.SAdd(key, value)
}

func SMembers(key string) ([]string, error) {
	members := redisDb.SMembers(key)
	result, err := members.Result()
	return result, err
}
