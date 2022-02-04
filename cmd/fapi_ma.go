package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/fapi/strategy"
)

func main() {
	const symbol = fapi.BCH
	fapi.RunStrategy(strategy.NewShow(symbol), symbol, "15m", 21)
}
