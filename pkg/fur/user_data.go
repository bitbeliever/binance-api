package fur

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/configs"
	"log"
	"time"
)

func UserDataStream(ch chan *futures.WsUserDataEvent) {
	c := futures.NewClient(configs.Cfg.Key.ApiKey, configs.Cfg.Key.SecretKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	listenKey, err := c.NewStartUserStreamService().Do(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("key:", listenKey)

	doneCh, stopCh, err := futures.WsUserDataServe(listenKey, func(event *futures.WsUserDataEvent) {
		ch <- event
	}, func(err error) {
		log.Println("data serve err", err)
		return
	})

	if err != nil {
		log.Println(err)
		return
	}
	_ = doneCh
	_ = stopCh

	go keepListenKeyAlive(c, listenKey)
}

func keepListenKeyAlive(client *futures.Client, listenKey string) {

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
