package main

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/helper"
	futures2 "github.com/bitbeliever/binance-api/pkg/fur"
	"log"
)

const (
	chBuf = 2 << 10 // (2 ^ 11)
)

func main() {

	//ws.KlineStream()

	log.Println(chBuf)
	ch := make(chan *futures.WsUserDataEvent, chBuf)
	futures2.UserDataStream(ch)

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
