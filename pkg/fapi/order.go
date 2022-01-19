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

根据 order type的不同，某些参数强制要求，具体如下:

Type	强制要求的参数
LIMIT	timeInForce, quantity, price
MARKET	quantity
STOP, TAKE_PROFIT	quantity, price, stopPrice
STOP_MARKET, TAKE_PROFIT_MARKET	stopPrice
TRAILING_STOP_MARKET	callbackRate
条件单的触发必须:

如果订单参数priceProtect为true:
达到触发价时，MARK_PRICE(标记价格)与CONTRACT_PRICE(合约最新价)之间的价差不能超过改symbol触发保护阈值
触发保护阈值请参考接口GET /fapi/v1/exchangeInfo 返回内容相应symbol中"triggerProtect"字段
STOP, STOP_MARKET 止损单:
买入: 最新合约价格/标记价格高于等于触发价stopPrice
卖出: 最新合约价格/标记价格低于等于触发价stopPrice
TAKE_PROFIT, TAKE_PROFIT_MARKET 止盈单:
买入: 最新合约价格/标记价格低于等于触发价stopPrice
卖出: 最新合约价格/标记价格高于等于触发价stopPrice
TRAILING_STOP_MARKET 跟踪止损单:
买入: 当合约价格/标记价格区间最低价格低于激活价格activationPrice,且最新合约价格/标记价高于等于最低价设定回调幅度。
卖出: 当合约价格/标记价格区间最高价格高于激活价格activationPrice,且最新合约价格/标记价低于等于最高价设定回调幅度。
TRAILING_STOP_MARKET 跟踪止损单如果遇到报错 {"code": -2021, "msg": "Order would immediately trigger."}
表示订单不满足以下条件:

买入: 指定的activationPrice 必须小于 latest price
卖出: 指定的activationPrice 必须大于 latest price
newOrderRespType 如果传 RESULT:

MARKET 订单将直接返回成交结果；
配合使用特殊 timeInForce 的 LIMIT 订单将直接返回成交或过期拒绝结果。
STOP_MARKET, TAKE_PROFIT_MARKET 配合 closePosition=true:

条件单触发依照上述条件单触发逻辑
条件触发后，平掉当时持有所有多头仓位(若为卖单)或当时持有所有空头仓位(若为买单)
不支持 quantity 参数
自带只平仓属性，不支持reduceOnly参数
双开模式下,LONG方向上不支持BUY; SHORT 方向上不支持SELL
{
    "clientOrderId": "testOrder", // 用户自定义的订单号
    "cumQty": "0",
    "cumQuote": "0", // 成交金额
    "executedQty": "0", // 成交量
    "orderId": 22542179, // 系统订单号
    "avgPrice": "0.00000",  // 平均成交价
    "origQty": "10", // 原始委托数量
    "price": "0", // 委托价格
    "reduceOnly": false, // 仅减仓
    "side": "SELL", // 买卖方向
    "positionSide": "SHORT", // 持仓方向
    "status": "NEW", // 订单状态
    "stopPrice": "0", // 触发价，对`TRAILING_STOP_MARKET`无效
    "closePosition": false,   // 是否条件全平仓
    "symbol": "BTCUSDT", // 交易对
    "timeInForce": "GTC", // 有效方法
    "type": "TRAILING_STOP_MARKET", // 订单类型
    "origType": "TRAILING_STOP_MARKET",  // 触发前订单类型
    "activatePrice": "9020", // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
    "priceRate": "0.3", // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段
    "updateTime": 1566818724722, // 更新时间
    "workingType": "CONTRACT_PRICE", // 条件价格触发类型
    "priceProtect": false            // 是否开启条件单触发保护
}
*/
/*
marginType 保证金模式: 全仓/逐仓 crossed/isolated
multiAssetsMargin 联合保证金模式: 单币种/跨币种
杠杆
positionSide 持仓方向: 单向持仓模式下非必填，默认且仅可填BOTH;在双向持仓模式下必填,且仅可选择 LONG 或 SHORT
*/
//func CreateOrder(symbol string, side futures.SideType, qty string) (*futures.CreateOrderResponse, error) {
//
//
//	client := NewClient()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
//	defer cancel()
//
//	order, err := client.NewCreateOrderService().
//		Symbol(symbol).
//		Side(side).
//		Type(futures.OrderTypeTakeProfitMarket).
//		Quantity(qty).
//		PositionSide(futures.PositionSideTypeBoth).    // 持仓方向 单向必填默认为BOTH
//		WorkingType(futures.WorkingTypeContractPrice). // stopPrice 触发类型: MARK_PRICE(标记价格), CONTRACT_PRICE(合约最新价). 默认 CONTRACT_PRICE
//		//StopPrice().                                   // 触发价 STOP, STOP_MARKET, TAKE_PROFIT, TAKE_PROFIT_MARKET 需要此参数
//		//Price("0.0030000"). // 委托价格
//		//closePosition(true). //true, false；触发后全部平仓，仅支持STOP_MARKET和TAKE_PROFIT_MARKET；不与quantity合用；自带只平仓效果，不与reduceOnly 合用
//		//PriceProtect() // 条件单触发保护："TRUE","FALSE", 默认"FALSE". 仅 STOP, STOP_MARKET, TAKE_PROFIT, TAKE_PROFIT_MARKET 需要此参数
//		Do(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	return order, nil
//}
func CreateOrder(symbol string, side futures.SideType, qty string) (*futures.CreateOrderResponse, error) {
	order, err := NewClient().NewCreateOrderService().
		Symbol(symbol).
		Side(side).
		Type(futures.OrderTypeMarket).
		Quantity(qty).
		PositionSide(futures.PositionSideTypeBoth).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	return order, nil
}

// CreateOrderDual 双向持仓
func CreateOrderDual(symbol string, side futures.SideType, positionSide futures.PositionSideType, qty string) (*futures.CreateOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	order, err := NewClient().NewCreateOrderService().
		Symbol(symbol).
		Side(side).
		Type(futures.OrderTypeMarket).
		Quantity(qty).
		PositionSide(positionSide).                    // 持仓方向 单向必填默认为BOTH
		WorkingType(futures.WorkingTypeContractPrice). // stopPrice 触发类型: MARK_PRICE(标记价格), CONTRACT_PRICE(合约最新价). 默认 CONTRACT_PRICE
		//StopPrice().                                   // 触发价 STOP, STOP_MARKET, TAKE_PROFIT, TAKE_PROFIT_MARKET 需要此参数
		//Price("0.0030000"). // 委托价格
		//closePosition(true). //true, false；触发后全部平仓，仅支持STOP_MARKET和TAKE_PROFIT_MARKET；不与quantity合用；自带只平仓效果，不与reduceOnly 合用
		//PriceProtect() // 条件单触发保护："TRUE","FALSE", 默认"FALSE". 仅 STOP, STOP_MARKET, TAKE_PROFIT, TAKE_PROFIT_MARKET 需要此参数
		Do(ctx)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// QueryOpenOrders /fapi/v1/openOrders 查询当前全部挂单
/*
  {
    "avgPrice": "0.00000",              // 平均成交价
    "clientOrderId": "abc",             // 用户自定义的订单号
    "cumQuote": "0",                        // 成交金额
    "executedQty": "0",                 // 成交量
    "orderId": 1917641,                 // 系统订单号
    "origQty": "0.40",                  // 原始委托数量
    "origType": "TRAILING_STOP_MARKET", // 触发前订单类型
    "price": "0",                   // 委托价格
    "reduceOnly": false,                // 是否仅减仓
    "side": "BUY",                      // 买卖方向
    "positionSide": "SHORT", // 持仓方向
    "status": "NEW",                    // 订单状态
    "stopPrice": "9300",                    // 触发价，对`TRAILING_STOP_MARKET`无效
    "closePosition": false,   // 是否条件全平仓
    "symbol": "BTCUSDT",                // 交易对
    "time": 1579276756075,              // 订单时间
    "timeInForce": "GTC",               // 有效方法
    "type": "TRAILING_STOP_MARKET",     // 订单类型
    "activatePrice": "9020", // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
    "priceRate": "0.3", // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段
    "updateTime": 1579276756075,        // 更新时间
    "workingType": "CONTRACT_PRICE", // 条件价格触发类型
    "priceProtect": false            // 是否开启条件单触发保护
  }
*/
func QueryOpenOrders(symbol string) {
	orders, err := NewClient().NewListOpenOrdersService().Symbol(symbol).Do(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("open orders", toJsonIndent(orders))
}

func QueryAllOpenOrders() ([]*futures.Order, error) {
	return NewClient().NewListOpenOrdersService().Symbol(ETH).Do(context.Background())

	//orders, err :=
	//if err != nil {
	//	log.Println(err)
	//	return orders, err
	//}
}

// QueryOrder 查询订单 /fapi/v1/order
func QueryOrder(symbol string, orderID int64) {
	order, err := NewClient().NewGetOrderService().Symbol(symbol).OrderID(orderID).Do(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(toJson(order))
}

// QueryAllOrders  查询所有订单(包括历史订单)  /fapi/v1/allOrders
func QueryAllOrders(symbol string) ([]*futures.Order, error) {
	orders, err := NewClient().NewListOrdersService().
		Limit(10).
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//log.Println(toJson(orders))
	log.Println("last order", toJson(orders[len(orders)-1]), len(orders), time.UnixMilli(orders[0].UpdateTime), time.UnixMilli(orders[len(orders)-1].UpdateTime).Format("15:04:05"))
	return orders, err
}

func CancelOrder(orderID int64) (*futures.CancelOrderResponse, error) {
	return NewClient().NewCancelOrderService().OrderID(orderID).Do(context.Background())
}
