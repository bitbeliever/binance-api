package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/spot"
	"time"
)

func main() {
	fapi.QueryOpenOrders()
	fapi.QueryAccountBalance()

	go fapi.UserDataStreamTest()
	//fapi.CreateOrder("BNBUSDT", futures.SideTypeBuy)

	time.Sleep(time.Second)
	fapi.QueryOpenOrders()
	fapi.QueryAccountBalance()
	fapi.QueryAccount()

	//fapi.QueryOrder("BNBUSDT", 37116496894)
	//fapi.QueryAllOrders(fapi.BNB)
	//spot.AccountService()
	spot.AccountService()
	fapi.QueryAllOrders(fapi.BNB)
	select {}
}
