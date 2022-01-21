package main

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

func main() {

	//helper.JsonLog(fapi.QueryBalance())
	//helper.JsonLog(fapi.QueryAccountAssets())
	helper.JsonLog(fapi.QueryAccountPositions())

	//closeOp()
	//fapi.CloseAllPositions()
	helper.JsonLog(fapi.QueryAccountPositions())
	//helper.JsonLog(fapi.QueryBalance())

	//fapi.ComparePNLTest()
}

func closeOp() {
	// Note: positionSideType不变, reverse sideType(buy/sell)
	_, err := fapi.CreateOrderDual(fapi.LTC, futures.SideTypeBuy, futures.PositionSideTypeShort, "0.070")
	if err != nil {
		log.Println(err)
	}
	//log.Println(helper.ToJsonIndent(resp))
}
