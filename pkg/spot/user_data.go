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

	doneCh, stopCh, err := binance.WsUserDataServe(listenKey, func(event *binance.WsUserDataEvent) {
		log.Println("spot event update")
	}, func(err error) {

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
	ticker := time.NewTicker(time.Minute * 30)
	for {
		select {
		case <-ticker.C:
			if err := client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(context.Background()); err != nil {
				log.Println(err)
			}
		}
	}
}
