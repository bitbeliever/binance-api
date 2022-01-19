package main

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

func main() {
	//helper.JsonLog(fapi.ModifyLeverage(fapi.BNB, 75))
	//helper.JsonLog(fapi.CreateOrder(fapi.LTC, futures.SideTypeSell, "0.05"))
	//helper.JsonLog(fapi.ModifyLeverage(fapi.ETH, 100))

	//closeOp()

	//helper.JsonLog(fapi.QueryBalance())
	//helper.JsonLog(fapi.QueryAccountAssets())
	helper.JsonLog(fapi.QueryAccountPositions())

	//fapi.CloseAllPositions()
	//closeOp()
	//
	//helper.JsonLog(fapi.QueryAccountPositions())
	helper.JsonLog(fapi.QueryAccountPositions())
}

func closeOp() {
	// Note: positionSideType不变, reverse sideType(buy/sell)
	_, err := fapi.CreateOrderDual(fapi.LTC, futures.SideTypeBuy, futures.PositionSideTypeShort, "0.070")
	if err != nil {
		log.Println(err)
	}
	//log.Println(helper.ToJsonIndent(resp))
}
