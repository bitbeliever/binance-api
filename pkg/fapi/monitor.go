package fapi

import "github.com/adshao/go-binance/v2/futures"

// temp
func monitor(symbol string, closePrice string, sub float64, side futures.SideType) {
	ch := AggTradePrice(symbol)

	for {
		select {
		case price := <-ch:
			// 当前价格高于下单价 平单
			if Str2Float64(closePrice)-Str2Float64(price) >= sub {
				CreateOrder(symbol, side, "0.05")
				return

			} else if Str2Float64(price)-Str2Float64(closePrice) >= sub {
				// 当前价格低于下单价价格
				CreateOrder(symbol, side, "0.05")
				return
			}
		}
	}
}
