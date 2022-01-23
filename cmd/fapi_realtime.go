package main

import (
	"github.com/bitbeliever/binance-api/configs"
	"github.com/bitbeliever/binance-api/pkg/account"
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/fapi/trade"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

func main() {
	p, err := account.QueryAccountPositions()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("positions", helper.ToJsonIndent(p))
	a, err := account.QueryAccountAssets()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("assets", helper.ToJson(a))

	const symbol = fapi.LTC

	// 全仓/逐仓设置
	mode, err := trade.PositionMode()
	if err != nil {
		log.Println(err)
		return
	}
	// 非双向持仓模式
	if !mode.DualSidePosition {
		if err := trade.PositionModeChange(true); err != nil {
			log.Println(err)
			return
		}
	}

	// 杠杆调整
	if err := trade.LeverageSetMax(symbol); err != nil {
		log.Println(err)
		return
	}

	fapi.RealTimeKline(symbol, configs.Cfg.KlineInterval)
}
