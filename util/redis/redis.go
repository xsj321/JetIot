package redis

import (
	"JetIot/conf"
	"github.com/go-redis/redis"
)

var redisDb *redis.Client

func init() {
	redisDb = redis.NewClient(&redis.Options{
		Addr:     conf.Default.RedisServer, // use default Addr
		Password: "",                       // no password set
		DB:       0,                        // use default DB
	})
}

func Set(key string, value interface{}) error {
	cmd := redisDb.Set(key, value, 0)
	return cmd.Err()
}

func Get(key string) (string, error) {
	get := redisDb.Get(key)
	return get.String(), get.Err()
}
