package fapi

import (
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"math"
)

/*
起始保证金= 成交数量x 开仓价格x IMR
	*初始保证金比率(IMR) = 1 / 杠杆

收益：
	做多= (平仓价格- 开仓价格) x 成交数量
	做空= (开仓价格- 平仓价格) x 成交数量

回报率(ROE) = 收益/ 起始保证金= 方向x (1 - 开仓价格/ 平仓价格) / IMR

目标价格：
	做多目标价格= 开仓价格* (回报率/ 杠杆+ 1)
	做空目标价格= 开仓价格* (1 - 回报率/ 杠杆)
=========

未实现盈亏(PNL) = 头寸大小* 订单方向* (最新价格- 开仓价格)
股本回报率(ROE%) = 以USDT计价的未实现盈亏(PNL) / 开仓保证金= ((最新价格- 开仓价格) * 订单方向* 规模) / (头寸金额* 合约倍数* 标记价格* 初始保证金比率(IMR))
订单方向：多头订单为1，空头订单为-1
*/

// 计算PNL: 未实现盈亏(PNL) = 头寸大小* 订单方向* (最新价格- 开仓价格)
func calcPNL(positionSize float64, sideType futures.PositionSideType, entry float64, price float64) float64 {
	var side float64
	if sideType == futures.PositionSideTypeLong {
		side = 1
	} else if sideType == futures.PositionSideTypeShort {
		side = -1
	} else {
		log.Println("wrong side type", sideType)
		return 0
	}

	return positionSize * side * (price - entry)
}

// 单个仓位, 监控, 止损pnl
// todo 通过 ORDER_TRADE_UPDATE 进行监控
func watchPNLStopLimit(order *futures.AccountPosition, stop chan struct{}) {
	ch, err := AggTradePrice(order.Symbol)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		select {
		case curPriceStr := <-ch:
			pnl := calcPNL(Str2Float64(order.PositionAmt), order.PositionSide, Str2Float64(order.EntryPrice), Str2Float64(curPriceStr))
			//log.Println("current pnl", pnl)
			// 触发止损
			if pnl < 0 && math.Abs(pnl) > principal.stopPNL() {
				log.Printf("触发止损 pnl: %v \t stopPNL %v \n", pnl, principal.stopPNL())
				if err := closePosition(order); err != nil {
					log.Println(err)
				}
				return
			}
		case <-stop:
			return
		}
	}
}

// pnl 对比测试
func ComparePNLTest() {
	pos, err := QueryAccountPositions()
	if err != nil {
		log.Println(err)
		return
	}

	p := pos[1]
	log.Println("position un_profit", p.UnrealizedProfit)
	//calcPNL(Str2Float64(p.PositionAmt), p.PositionSide, Str2Float64(p.EntryPrice), ch)
}
