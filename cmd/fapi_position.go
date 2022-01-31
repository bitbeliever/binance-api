package main

import (
	"github.com/bitbeliever/binance-api/pkg/account"
	"github.com/bitbeliever/binance-api/pkg/cache"
	"github.com/bitbeliever/binance-api/pkg/fapi/position"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

func main() {

	//helper.JsonLog(fapi.QueryBalance())
	//helper.JsonLog(fapi.QueryAccountAssets())
	helper.JsonLog(account.QueryAccountPositions())

	//closeOp()
	closeAndClearKeys()
	helper.JsonLog(account.QueryAccountPositions())
	helper.JsonLog(account.QueryBalance())

	//fapi.ComparePNLTest()
}

func closeAndClearKeys() {
	position.CloseAllPositions()
	if err := cache.ClearKeys("smooth_*"); err != nil {
		log.Println(err)
	}

}
