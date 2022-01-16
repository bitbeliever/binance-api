package spot

import (
	"github.com/adshao/go-binance/v2"
	"log"
)

const (
	chBuf = 2 << 10
)

func KlineStream() {
	//c := client.NewClient()
	klineCh := make(chan binance.WsKline, chBuf)
	bollCh := make(chan bollResult, chBuf)

	go func() {
		doneCh, stopCh, err := binance.WsKlineServe(LTC, "15m", func(event *binance.WsKlineEvent) {
			//log.Println(toJson(event))
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

	go subscribeWsKline(klineCh, bollCh)
	for {
		select {
		case bRes := <-bollCh:
			//log.Println(bRes)
			log.Println("spot", toJson(bRes))
		}
	}

}
