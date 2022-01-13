package main

import (
	"github.com/adshao/go-binance/v2/futures"
	futures2 "github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"github.com/bitbeliever/binance-api/pkg/spot"
	"log"
)

const (
	chBuf = 2 << 10 // (2 ^ 11)
)

func main() {
	//ws.KlineStream()
	ch := make(chan *futures.WsUserDataEvent, chBuf)
	futures2.UserDataStream(ch)
	spot.UserDataStream()

	for {
		select {
		case event := <-ch:
			log.Println(helper.ToJson(event))
			if event.Event == futures.UserDataEventTypeAccountUpdate { // ACCOUNT_UPDATE
				log.Println("futures user data")

			}

		}
	}
}
