package fapi

import (
	"github.com/adshao/go-binance/v2/futures"
	"log"
)

const (
	chBuf = 2 << 10
)

func KlineStream(symbol, interval string) chan futures.WsKline {
	klineCh := make(chan futures.WsKline, chBuf)
	//bollCh := make(chan bollResult, chBuf)

	go func() {
		doneCh, stopCh, err := futures.WsKlineServe(symbol, interval, func(event *futures.WsKlineEvent) {
			//k := event.Kline
			//log.Println(time.UnixMilli(event.Time).Format("15:04:05"), time.UnixMilli(k.StartTime).Format("15:04:05"), time.UnixMilli(k.EndTime).Format("15:04:05"), k.Open, k.High, k.Low, k.Close)
			klineCh <- event.Kline

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
	//		_ = bRes
	//		//log.Println(bRes)
	//		//log.Println("future", toJson(bRes))
	//	}
	//}

	return klineCh
}
