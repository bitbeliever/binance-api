package strategy

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/fapi/indicator"
	"github.com/bitbeliever/binance-api/pkg/fapi/order"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

/*
技术指标选择平均价。
以比特币举例，当比特币30分钟线的实际价格和average price相差达到200美金，我们开单，在水位线上我们开空单，水位线以下我们开空单，平仓价以average price 的动态价格为准。
如果价差相差超过300美金，我们可以继续加仓。不用设置止损。因为始终会回归到average price的价格。这个就更加简化了我们操作流程。而且全币种都可以同时做。
我刚才用30分钟线测了比特币和莱特币，都可以用这种方式，莱特币设置开仓价，是在average price和实际价格价差达到5毛开仓，达到一块继续加仓。
你看看用这种是不是更方便点，和布林通道逻辑一样

可以同时做所有币种，只要平均价和现价有足够的价差，我们就可以从中套利。现在需要确定的是，这中间的利润，能否覆盖我们的手续费和利息。不设置止损，是因为平均价是动态的，只要平均价是我们的平仓价，就随时都可以平仓。需要的就是价格的波动。

*/

type Average struct {
	symbol string
	gap    float64
	short  *futures.CreateOrderResponse
	long   *futures.CreateOrderResponse
}

func NewAverage(symbol string, gap float64) *Average {
	return &Average{
		symbol: symbol,
		gap:    gap,
	}
}

func (a *Average) Do(lines []*futures.Kline) error {
	ma := indicator.Ind(lines).Ma()
	boll := indicator.Ind(lines).Boll()
	if !boll.IsInsideBand() {
		return nil
	}

	// 位于ma线上方 sell short
	if helper.Str2Float64(ma.CurrentPrice()) >= helper.Str2Float64(ma.AveragePrice())+a.gap {
		if a.short == nil {
			resp, err := order.DualSellShortSL(a.symbol, "0.08", ma.AveragePrice())
			if err != nil {
				return err
			}
			log.Printf("avg sell short avg: %s, current: %s order %v  \n", ma.AveragePrice(), ma.CurrentPrice(), helper.ToJson(resp))
			if err != nil {
				return err
			}
			a.short = resp
		}
	} else if helper.Str2Float64(ma.CurrentPrice()) <= helper.Str2Float64(ma.AveragePrice())-a.gap { // 位于ma线下方  buy long
		if a.long == nil {
			resp, err := order.DualBuyLongSL(a.symbol, "0.08", ma.AveragePrice())
			if err != nil {
				return err
			}
			log.Printf("avg buy long avg: %s, current: %s order %v  \n", ma.AveragePrice(), ma.CurrentPrice(), helper.ToJson(resp))
			a.long = resp
		}
	}

	return nil
}

func (a *Average) resume() {

}
