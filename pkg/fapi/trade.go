package fapi

import (
	"context"
	"fmt"
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

// LeverageBracket 杠杆分层标准 /fapi/v1/leverageBracket
func LeverageBracket(symbol string) ([]*futures.LeverageBracket, error) {
	return NewClient().NewGetLeverageBracketService().Symbol(symbol).Do(context.Background())
}

// LeverageSetMax 设置该交易对最大杠杆 todo 杠杆分层
func LeverageSetMax(symbol string) error {
	brackets, err := LeverageBracket(symbol)
	if err != nil {
		return err
	}
	if len(brackets) == 0 {
		return fmt.Errorf("error brackets data %v", brackets)
	}
	if brackets[0].Symbol != symbol {
		return fmt.Errorf("symbol incorrespond %v %v", symbol, brackets[0].Symbol)
	}

	leverage, err := NewClient().NewChangeLeverageService().Symbol(symbol).Leverage(brackets[0].Brackets[0].InitialLeverage).Do(context.Background())
	if err != nil {
		return err
	}
	log.Println("杠杆设置最大:", toJson(leverage))

	return nil
}

// PositionMode 查询持仓模式 "true": 双向持仓模式；"false": 单向持仓模式
func PositionMode() (*futures.PositionMode, error) {
	return NewClient().NewGetPositionModeService().Do(context.Background())
}

// PositionModeChange 更改持仓模式
func PositionModeChange(dual bool) error {
	return NewClient().NewChangePositionModeService().DualSide(dual).Do(context.Background())
}

func getLatestPrice(symbol string) {

}