package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/fapi/order"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

func main() {
	resp, err := order.DualBuyLong(fapi.LTC, "0.1")
	if err != nil {
		log.Println(err)
	}
	log.Println(helper.ToJson(resp))
}
