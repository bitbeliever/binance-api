package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"log"
)

func main() {
	log.Println(fapi.QueryAllOpenOrders())
	fapi.QueryAccountBalance()

	go fapi.RecvUserDataStream()
	//fapi.CreateOrder("BNBUSDT", futures.SideTypeBuy)

	fapi.QueryAccountBalance()
	fapi.QueryAccount()
	select {}
}
