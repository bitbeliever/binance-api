package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"log"
)

func main() {
	log.Println(fapi.QueryAllOpenOrders())
	log.Println(fapi.QueryBalance())

	go fapi.RecvUserDataStream()

	log.Println(fapi.QueryBalance())
	//fapi.QueryAccount()
	//fapi.QueryAllOrders(fapi.BNB)
	select {}
}
