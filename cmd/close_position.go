package main

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/fapi"
)

func main() {
	// sell position
	fapi.CreateOrder(fapi.BNB, futures.SideTypeSell)
}
