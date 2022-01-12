package fapi

import (
	"bytes"
	"encoding/json"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"time"
)

type candle struct {
	StartTime int64  `json:"t"`
	EndTime   int64  `json:"T"`
	Symbol    string `json:"s"`
	Open      string `json:"o"`
	Close     string `json:"c"`
	High      string `json:"h"`
	Low       string `json:"l"`
}

const (
	format = "15:04:05"
)

// KlineStreamTest K线stream逐秒推送所请求的K线种类(最新一根K线)的更新
func KlineStreamTest() {

	go func() {
		doneCh, stopCh, err := futures.WsKlineServe("ETHUSDT", "15m", func(event *futures.WsKlineEvent) {
			//log.Println(toJson(event))
			log.Println("klinetest time", time.UnixMilli(event.Time).Format(format), event.Kline.Open, event.Kline.High, event.Kline.Low, event.Kline.Close)

		}, func(err error) {
			log.Println(err)
		})

		_ = doneCh
		_ = stopCh
		if err != nil {
			log.Println(err)
			return
		}
	}()

	//go subscribeWsKline(klineCh, bollCh)
	//for {
	//	select {
	//	case bRes := <-bollCh:
	//		//log.Println(bRes)
	//		log.Println("future", toJson(bRes))
	//	}
	//}

	select {}
}

func encode(v interface{}) interface{} {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	var candle candle
	if err := json.NewDecoder(bytes.NewBuffer(b)).Decode(&candle); err != nil {
		panic(err)
	}

	return candle
}
