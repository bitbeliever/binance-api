package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/fapi/strategy"
)

func main() {
	const symbol = fapi.LTC
	fapi.RunStrategy(strategy.NewShow(), symbol, "15m", 21)
}