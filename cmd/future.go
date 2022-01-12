package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
)

func main() {
	go fapi.KlineStream("ETHUSDT", "15m")
	//go spotws.KlineStream()

	//fapi.AccountService()
	select {}
}
