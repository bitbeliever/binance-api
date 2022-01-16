package fapi

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"log"
)

// AggTrade 最新合约价格
func AggTrade(symbol string) chan *futures.WsAggTradeEvent {
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
		log.Println(err)
	}

	return ch
}

func AggTradePrice(symbol string) chan string {
	ch := make(chan string, 512)
	_, _, err := futures.WsAggTradeServe(symbol, func(event *futures.WsAggTradeEvent) {
		ch <- event.Price
	}, func(err error) {
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Println(err)
	}

	return ch
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

// LeverageBrackets Leverage 杠杆分层标准
func LeverageBrackets() {
	lb, err := NewClient().NewGetLeverageBracketService().Symbol(ETH).Do(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(toJson(lb))
}

// ModifyLeverage 调制开仓杠杆   "maxNotionalValue": "1000000", // 当前杠杆倍数下允许的最大名义价值
func ModifyLeverage(symbol string, leverage int) (*futures.SymbolLeverage, error) {
	res, err := NewClient().NewChangeLeverageService().Symbol(symbol).Leverage(leverage).Do(context.Background())
	if err != nil {
		log.Println(err)
		return res, nil
	}

	return res, nil
}

// PositionMode 查询持仓模式 "true": 双向持仓模式；"false": 单向持仓模式
func PositionMode() {
	mode, err := NewClient().NewGetPositionModeService().Do(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(toJson(mode))
}

// PositionModeChange 更改持仓模式
func PositionModeChange(dual bool) {
	err := NewClient().NewChangePositionModeService().DualSide(dual).Do(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
}
