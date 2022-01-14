package fapi

import (
	"github.com/adshao/go-binance/v2/futures"
	"log"
)

// AggTrade 最新合约价格
func AggTrade(symbol string) chan *futures.WsAggTradeEvent {
	ch := make(chan *futures.WsAggTradeEvent, 512)
	_, _, err := futures.WsAggTradeServe(symbol, func(event *futures.WsAggTradeEvent) {
		//log.Println(time.UnixMilli(event.TradeTime).Format("15:04:05"), toJson(event))
		ch <- event
	}, func(err error) {
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Println(err)
	}

	return ch
}

func AggTradePrice(symbol string) chan string {
	ch := make(chan string, 512)
	_, _, err := futures.WsAggTradeServe(symbol, func(event *futures.WsAggTradeEvent) {
		ch <- event.Price
	}, func(err error) {
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Println(err)
	}

	return ch
}