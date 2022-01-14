package main

import "github.com/bitbeliever/binance-api/pkg/fapi"

func main() {
	fapi.RealTimeKline(fapi.ETH, "15m")
}
