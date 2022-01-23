package main

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/account"
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/fapi/order"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

func main() {

	//helper.JsonLog(fapi.QueryBalance())
	//helper.JsonLog(fapi.QueryAccountAssets())
	helper.JsonLog(account.QueryAccountPositions())

	//closeOp()
	//fapi.CloseAllPositions()
	helper.JsonLog(account.QueryAccountPositions())
	helper.JsonLog(account.QueryBalance())

	//fapi.ComparePNLTest()
}

func closeOp() {
	// Note: positionSideType不变, reverse sideType(buy/sell)
	_, err := order.CreateOrderDual(fapi.LTC, futures.SideTypeBuy, futures.PositionSideTypeShort, "0.070")
	if err != nil {
		log.Println(err)
	}
	//log.Println(helper.ToJsonIndent(resp))
}
