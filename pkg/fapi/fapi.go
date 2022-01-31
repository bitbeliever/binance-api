package fapi

import (
	"github.com/bitbeliever/binance-api/pkg/account"
	"github.com/bitbeliever/binance-api/pkg/fapi/strategy"
	"github.com/bitbeliever/binance-api/pkg/fapi/trade"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

func InitSettings(symbol string) {

	p, err := account.QueryAccountPositions()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("positions", helper.ToJson(p))
	a, err := account.QueryAccountAssets()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("assets", helper.ToJson(a))

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
	lev, err := trade.LeverageSetMax(symbol)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("杠杆设置:", lev.Leverage)
}

func RunStrategy(s strategy.Strategy, symbol string, interval string, limit int) {
	InitSettings(symbol)
	ch, err := RealTimeKline(symbol, interval, limit)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		select {
		case lines := <-ch:
			err := s.Do(lines)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
