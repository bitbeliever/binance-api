package fapi

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"log"
)

func init() {
}

// AggTrade 最新合约价格
func AggTrade(symbol string) (chan *futures.WsAggTradeEvent, error) {
	ch := make(chan *futures.WsAggTradeEvent, 512)
	_, _, err := futures.WsAggTradeServe(symbol, func(event *futures.WsAggTradeEvent) {
		//log.Println(time.UnixMilli(event.TradeTime).Format("15:04:05"), toJson(event))
		ch <- event
	}, func(err error) {
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func AggTradePrice(symbol string) (chan string, error) {
	ch := make(chan string, 512)
	_, _, err := futures.WsAggTradeServe(symbol, func(event *futures.WsAggTradeEvent) {
		ch <- event.Price
	}, func(err error) {
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return nil, err
	}

	return ch, nil
}

// UpdateMarginType 变换全逐仓模式 isolated || crossed  /fapi/v1/marginType
func UpdateMarginType(symbol string, typ futures.MarginType) error {
	return NewClient().NewChangeMarginTypeService().Symbol(symbol).MarginType(typ).Do(context.Background())
}

// ModifyIsolatedMargin 调整逐仓保证金 /fapi/v1/positionMargin
func ModifyIsolatedMargin(symbol string) {
	// todo
	//NewClient().NewUpdatePositionMarginService().Symbol(symb)
}

// ModifyLeverage 调制开仓杠杆   "maxNotionalValue": "1000000", // 当前杠杆倍数下允许的最大名义价值
func ModifyLeverage(symbol string, leverage int) (*futures.SymbolLeverage, error) {
	return NewClient().NewChangeLeverageService().Symbol(symbol).Leverage(leverage).Do(context.Background())
}

// PositionMode 查询持仓模式 "true": 双向持仓模式；"false": 单向持仓模式
func PositionMode() (*futures.PositionMode, error) {
	return NewClient().NewGetPositionModeService().Do(context.Background())
}

// PositionModeChange 更改持仓模式
func PositionModeChange(dual bool) error {
	return NewClient().NewChangePositionModeService().DualSide(dual).Do(context.Background())
}
