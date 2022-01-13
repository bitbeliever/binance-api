package fapi

import (
	"github.com/adshao/go-binance/v2/futures"
	"log"
)

func WSTicker() {
	go func() {
		_, _, err := futures.WsMiniMarketTickerServe(LTC, func(event *futures.WsMiniMarketTickerEvent) {
			log.Println(toJson(event))
		}, func(err error) {
			if err != nil {
				log.Println(err)
			}
		})

		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		_, _, err := futures.WsBookTickerServe(LTC, func(event *futures.WsBookTickerEvent) {
			log.Println(toJson(event))

		}, func(err error) {
			if err != nil {
				log.Println(err)
			}
		})

		if err != nil {
			log.Println(err)
			return
		}
	}()

	select {}
}
