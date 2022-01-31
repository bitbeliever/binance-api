package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/fapi/strategy"
)

func main() {
	const symbol = fapi.DOG
	fapi.RunStrategy(strategy.NewAverage(symbol, 0.1), symbol, "30m", 7)
}
