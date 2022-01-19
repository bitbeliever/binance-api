package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/helper"
)

func main() {
	helper.JsonLog(fapi.QueryAccountPositions())
	//fapi.QueryAccount()
	//fapi.QueryAllOrders(fapi.BNB)
}
