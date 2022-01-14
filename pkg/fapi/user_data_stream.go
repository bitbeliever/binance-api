package fapi

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/configs"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
	"time"
)

func UserDataStream(ch chan *futures.WsUserDataEvent) {
	c := futures.NewClient(configs.Cfg.Key.ApiKey, configs.Cfg.Key.SecretKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	listenKey, err := c.NewStartUserStreamService().Do(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("fapi key:", listenKey)

	doneCh, stopCh, err := futures.WsUserDataServe(listenKey, func(event *futures.WsUserDataEvent) {
		ch <- event
	}, func(err error) {
		log.Println("fapi data serve err", err)
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

// keep alive
func keepListenKeyAlive(client *futures.Client, listenKey string) {

	ticker := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-ticker.C:
			log.Println("user data stream keeping alive")
			if err := client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(context.Background()); err != nil {
				log.Println(err)
			}
		}
	}
}

func UserDataStreamTest() {
	//ws.KlineStream()
	ch := make(chan *futures.WsUserDataEvent, chBuf)
	UserDataStream(ch)

	for {
		select {
		case event := <-ch:
			log.Println(helper.ToJson(event))

			//if event.Event == futures.UserDataEventTypeAccountUpdate { // ACCOUNT_UPDATE
			//	log.Println("futures user data")
			//}

		}
	}
}
