package fapi

import (
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"time"
)

// AggTrade 最新合约价格
func AggTrade() {
	_, _, err := futures.WsAggTradeServe(LTC, func(event *futures.WsAggTradeEvent) {
		log.Println(time.UnixMilli(event.TradeTime).Format("15:04:05"), toJson(event))
	}, func(err error) {
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Println(err)
	}

	select {}
}
