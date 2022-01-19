package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

func main() {
	p, err := fapi.QueryAccountPositions()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("positions", helper.ToJsonIndent(p))
	a, err := fapi.QueryAccountAssets()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("assets", helper.ToJson(a))

	const symbol = fapi.LTC
	log.Println("doing symbol:", symbol)

	// 全仓/逐仓设置
	mode, err := fapi.PositionMode()
	if err != nil {
		log.Println(err)
		return
	}
	// 非双向持仓模式
	if !mode.DualSidePosition {
		if err := fapi.PositionModeChange(true); err != nil {
			log.Println(err)
			return
		}
	}

	// 杠杆调整
	if err := fapi.LeverageSetMax(symbol); err != nil {
		log.Println(err)
		return
	}

	fapi.RealTimeKline(symbol, "15m")
}
