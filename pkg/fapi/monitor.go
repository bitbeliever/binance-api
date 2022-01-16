package fapi

import (
	"github.com/adshao/go-binance/v2/futures"
	"log"
)

// temp
func monitor(symbol string, closePriceStr string, sub float64, side futures.SideType) {
	ch := AggTradePrice(symbol)

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
