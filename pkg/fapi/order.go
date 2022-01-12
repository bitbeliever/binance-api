package fapi

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"time"
)

/*
CreateOrder

symbol	STRING	YES	交易对
side	ENUM	YES	买卖方向 SELL, BUY
positionSide	ENUM	NO	持仓方向，单向持仓模式下非必填，默认且仅可填BOTH;在双向持仓模式下必填,且仅可选择 LONG 或 SHORT
type	ENUM	YES	订单类型 LIMIT, MARKET, STOP, TAKE_PROFIT, STOP_MARKET, TAKE_PROFIT_MARKET, TRAILING_STOP_MARKET
reduceOnly	STRING	NO	true, false; 非双开模式下默认false；双开模式下不接受此参数； 使用closePosition不支持此参数。
quantity	DECIMAL	NO	下单数量,使用closePosition不支持此参数。
price	DECIMAL	NO	委托价格
newClientOrderId	STRING	NO	用户自定义的订单号，不可以重复出现在挂单中。如空缺系统会自动赋值。必须满足正则规则 ^[\.A-Z\:/a-z0-9_-]{1,36}$
stopPrice	DECIMAL	NO	触发价, 仅 STOP, STOP_MARKET, TAKE_PROFIT, TAKE_PROFIT_MARKET 需要此参数
closePosition	STRING	NO	true, false；触发后全部平仓，仅支持STOP_MARKET和TAKE_PROFIT_MARKET；不与quantity合用；自带只平仓效果，不与reduceOnly 合用
activationPrice	DECIMAL	NO	追踪止损激活价格，仅TRAILING_STOP_MARKET 需要此参数, 默认为下单当前市场价格(支持不同workingType)
callbackRate	DECIMAL	NO	追踪止损回调比例，可取值范围[0.1, 5],其中 1代表1% ,仅TRAILING_STOP_MARKET 需要此参数
timeInForce	ENUM	NO	有效方法
workingType	ENUM	NO	stopPrice 触发类型: MARK_PRICE(标记价格), CONTRACT_PRICE(合约最新价). 默认 CONTRACT_PRICE
priceProtect	STRING	NO	条件单触发保护："TRUE","FALSE", 默认"FALSE". 仅 STOP, STOP_MARKET, TAKE_PROFIT, TAKE_PROFIT_MARKET 需要此参数
newOrderRespType	ENUM	NO	"ACK", "RESULT", 默认 "ACK"
recvWindow	LONG	NO
timestamp	LONG	YES
*/
func CreateOrder(symbol string, side futures.SideType) {
	client := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	order, err := client.NewCreateOrderService().Symbol("BNBETH").
		Side(futures.SideTypeBuy).Type(futures.OrderTypeLimit).
		TimeInForce(futures.TimeInForceTypeGTC).Quantity("5").
		Price("0.0030000").Do(ctx)
	if err != nil {
		log.Println(err)
	}

	log.Println(toJson(order))
}
