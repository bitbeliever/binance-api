package spot

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"github.com/bitbeliever/binance-api/configs"
	"log"
	"time"
)

func UserDataStream() {
	c := binance.NewClient(configs.Cfg.Key.ApiKey, configs.Cfg.Key.SecretKey)
	listenKey, err := c.NewStartUserStreamService().Do(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("spot listen key", listenKey)

	doneCh, stopCh, err := binance.WsUserDataServe(listenKey, func(event *binance.WsUserDataEvent) {
		log.Println("spot event update")
		log.Println(toJson(event))
	}, func(err error) {
		log.Println("user data service err", err)

	})
	if err != nil {
		log.Println(err)
		return
	}

	_ = doneCh
	_ = stopCh
	go keepUserDataServiceAlive(c, listenKey)

}

func keepUserDataServiceAlive(client *binance.Client, listenKey string) {
	ticker := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-ticker.C:
			if err := client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(context.Background()); err != nil {
				log.Println(err)
			}
		}
	}
}
