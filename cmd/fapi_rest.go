package main

import (
	"github.com/bitbeliever/binance-api/pkg/account"
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/fapi/order"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

func main() {
	helper.JsonLog(account.QueryAccountPositions())
	//fapi.QueryAccount()
	//fapi.QueryAllOrders(fapi.BNB)
	o, err := order.DualBuyLong(fapi.BCH, "0.2")
	if err != nil {
		log.Println(err)
	}
	log.Println(helper.ToJson(o))
}
