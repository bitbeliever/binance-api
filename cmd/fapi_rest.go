package main

import (
	"github.com/bitbeliever/binance-api/pkg/account"
	"github.com/bitbeliever/binance-api/pkg/helper"
)

func main() {
	helper.JsonLog(account.QueryAccountPositions())
	//fapi.QueryAccount()
	//fapi.QueryAllOrders(fapi.BNB)
}
