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
func pnlMonitor(entry float64, qty string, ch chan float64) float64 {
	//entry := Str2Float64(position.EntryPrice)
	var pnl float64
	_ = pnl
	for {
		select {
		case curPrice := <-ch:
			pnl = (curPrice - entry) * Str2Float64(qty)
		}
	}

	return 0
}

// 单个仓位pnl
// 未实现盈亏(PNL) = 头寸大小* 订单方向* (最新价格- 开仓价格)
// todo positionSize == positionAmt?
func pnl(positionSize float64, sideType futures.PositionSideType, entry float64, ch chan float64) {
	var side float64
	if sideType == futures.PositionSideTypeLong {
		side = 1
	} else if sideType == futures.PositionSideTypeShort {
		side = -1
	} else {
		log.Println("wrong side type", sideType)
		return
	}

	for {
		select {
		case curPrice := <-ch:
			pnl := positionSize * side * (curPrice - entry)
			if pnl < 0 && math.Abs(pnl) > principal.stopBalance() {
				// todo stop

			}
		}
	}

}
