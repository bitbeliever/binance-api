package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/fapi/order"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

func main() {

	//if err := order.CancelAllOpenOrders(); err != nil {
	//	log.Println(err)
	//	return
	//}

	orders, err := order.QueryOpenOrders(fapi.LTC)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("open orders", helper.ToJson(orders))

	orders, err = order.QueryAllOpenOrders()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("all", helper.ToJson(orders), len(orders))

	//o, err := order.DualSellShortSL(fapi.LTC, "0.05", "124")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//log.Println(helper.ToJson(o))

}
