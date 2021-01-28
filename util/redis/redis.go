package redis

import (
	"JetIot/conf"
	"github.com/go-redis/redis"
)

var redisDb *redis.Client

func InitRedis() {
	redisDb = redis.NewClient(&redis.Options{
		Addr:     conf.Default.RedisServer + ":" + conf.Default.RedisPort, // use default Addr
		Password: "",                                                      // no password set
		DB:       0,                                                       // use default DB
	})
}

func Set(key string, value interface{}) error {
	err := redisDb.Set(key, value, 0).Err()
	return err
}

func Get(key string) ([]byte, error) {
	get := redisDb.Get(key)
	return get.Bytes()
}
