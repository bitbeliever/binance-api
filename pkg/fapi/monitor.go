package fapi

import (
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"math"
)

// temp
func monitor(symbol string, closePriceStr string, sub float64, side futures.SideType) {
	ch, err := AggTradePrice(symbol)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		select {
		case priceStr := <-ch:
			price := Str2Float64(priceStr)
			closePrice := Str2Float64(closePriceStr)

			// buy
			if side == futures.SideTypeBuy {
				// 当前价格高于下单价 平单
				if price-closePrice >= sub {
					CreateOrder(symbol, futures.SideTypeSell, "0.05")
					return
				} else if closePrice-price >= sub {
					// 当前价格低于下单价价格
					CreateOrder(symbol, futures.SideTypeSell, "0.05")
					return
				}
			} else if side == futures.SideTypeSell { // sell
				if closePrice-price >= sub {
					CreateOrder(symbol, futures.SideTypeBuy, "0.05")
					return
				} else if price-closePrice >= sub { // buy back
					CreateOrder(symbol, futures.SideTypeBuy, "0.05")
					return
				}
			} else {
				log.Println("Wrong SideType")
			}
		}
	}
}

func (s *doubleOpenStrategy) monitorOrderTP(chUpper, chLower chan struct{}) {
	for {
		select {
		case <-chUpper:
			if s.longOrder != nil {
				log.Println("触发多单止盈")
				if err := closePositionByOrderResp(s.longOrder); err != nil {
					log.Println(err)
					return
				}
				s.longOrder = nil
				return
			}

		case <-chLower:
			if s.shortOrder != nil {
				if err := closePositionByOrderResp(s.shortOrder); err != nil {
					log.Println(err)
					return
				}
				log.Println("触发空单止盈")
				s.shortOrder = nil
				return
			}
		}

	}
}

func (s *doubleOpenStrategy) monitorOrderSL(p *futures.CreateOrderResponse) {
	ch, err := AggTradePrice(p.Symbol)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		select {
		case curPriceStr := <-ch:
			pnl := calcPNL(Str2Float64(p.OrigQuantity), p.PositionSide, Str2Float64(p.AvgPrice), Str2Float64(curPriceStr))
			//log.Println("current pnl", pnl)
			// 触发止损
			if pnl < 0 && math.Abs(pnl) > principal.stopPNL() {
				log.Printf("触发止损 pnl: %v \t stopPNL %v \n", pnl, principal.stopPNL())
				if err := closePositionByOrderResp(p); err != nil {
					log.Println(err)
				}
				return
			}
			//case <-stop:
			//	return
		}
	}
}
