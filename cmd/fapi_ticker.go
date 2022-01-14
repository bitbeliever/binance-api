package main

import "github.com/bitbeliever/binance-api/pkg/fapi"

func main() {
	//go fapi.WSTicker()
	fapi.AggTrade(fapi.ETH)
}
