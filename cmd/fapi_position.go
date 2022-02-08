package main

import (
	"github.com/bitbeliever/binance-api/pkg/account"
	"github.com/bitbeliever/binance-api/pkg/fapi/position"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

func main() {

	helper.JsonLog(account.QueryAccountPositions())

	//_, err := order.DualSellShort("LTCUSDT", "0.5")
	//if err != nil {
	//	log.Println(err)
	//}
	//f, err := position.CloseAllPositionsBySymbol("BCHUSDT")
	//log.Println(f, err)
	//account.PositionsFormat()
	//closeAndClearKeys()
	helper.JsonLog(account.QueryBalance())

	//fapi.ComparePNLTest()
}

func closeAndClearKeys() {
	profit, err := position.CloseAllPositions()
	if err != nil {
		log.Println(err)
	}
	log.Println("close all total profit", profit)

	//if err := cache.ClearKeys("smooth_*"); err != nil {
	//	log.Println(err)
	//}
}
