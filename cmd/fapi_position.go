package main

import (
	"github.com/bitbeliever/binance-api/pkg/account"
	"github.com/bitbeliever/binance-api/pkg/helper"
)

func main() {

	//helper.JsonLog(fapi.QueryBalance())
	//helper.JsonLog(fapi.QueryAccountAssets())
	helper.JsonLog(account.QueryAccountPositions())

	//closeOp()
	//position.CloseAllPositions()
	helper.JsonLog(account.QueryAccountPositions())
	helper.JsonLog(account.QueryBalance())

	//fapi.ComparePNLTest()
}

func open() {
	//order.DualBuyLong(fapi.LTC)
}
