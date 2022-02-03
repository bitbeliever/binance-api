package main

import (
	"github.com/bitbeliever/binance-api/configs"
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/fapi/strategy"
)

func main() {
	//const symbol = fapi.LTC
	var symbol = configs.Cfg.Symbol
	fapi.RunStrategy(strategy.NewSmooth(symbol), symbol, "15m", 21)
	//fapi.RunStrategy(strategy.NewShow(), symbol, "15m", 21)
	//fapi.RunStrategy(strategy.NewSmoothTP(symbol), symbol, "15m", 21)
}
