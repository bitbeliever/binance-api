package main

import "github.com/bitbeliever/binance-api/pkg/fapi"

func main() {
	fapi.QueryOpenOrders()
	fapi.QueryAccountBalance()
}
