package spot

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/bitbeliever/binance-api/configs"
	"log"
	"time"
)

var (
	apiKey    = configs.Cfg.Key.ApiKey
	secretKey = configs.Cfg.Key.SecretKey
)

func Test() {
	client := binance.NewClient(apiKey, secretKey)
	log.Println(client.BaseURL)
	//futuresClient := binance.NewFuturesClient(apiKey, secretKey)   // USDT-M Futures
	//deliveryClient := binance.NewDeliveryClient(apiKey, secretKey) // Coin-M Futures
	klines, err := client.NewKlinesService().Symbol("LTCBTC").
		Interval("15m").Limit(1000).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, k := range klines {
		log.Println(time.UnixMilli(k.OpenTime), time.UnixMilli(k.CloseTime), toJson(k))
	}
}
