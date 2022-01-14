package main

import "github.com/bitbeliever/binance-api/pkg/fapi"

func main() {
	fapi.RealTimeKline(fapi.BNB, "15m")
}
