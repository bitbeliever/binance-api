package cache

import (
	"fmt"
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
		Addr:     fmt.Sprintf("%s:%d", configs.Cfg.Redis.Host, configs.Cfg.Redis.Port),
		Password: configs.Cfg.Redis.Password,
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
