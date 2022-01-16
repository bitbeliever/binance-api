package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"log"
)

func main() {
	log.Println(fapi.QueryAllOpenOrders())
	fapi.QueryAccountBalance()

	go fapi.RecvUserDataStream()

	fapi.QueryAccountBalance()
	fapi.QueryAccount()
	//fapi.QueryAllOrders(fapi.BNB)
	select {}
}
