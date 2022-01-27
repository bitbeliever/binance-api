package cache

import (
	"github.com/bitbeliever/binance-api/configs"
	"github.com/go-redis/redis"
	"log"
)

var Client *redis.Client

const (
	KeySmooth = "smooth_"
)

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr: configs.Cfg.Redis.Host,
	})

	res, err := Client.Ping().Result()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println(res)
}

func ClearKeys(pattern string) error {
	keys, err := Client.Keys(pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}

	return Client.Del(keys...).Err()
}
