package fapi

/*
## 合约状态 (contractStatus, status):
	PENDING_TRADING 待上市
	TRADING 交易中
	PRE_DELIVERING 预交割
	DELIVERING 交割中
	DELIVERED 已交割
	PRE_SETTLE 预结算
	SETTLING 结算中
	CLOSE 已下架


## 订单状态 (status):
	NEW 新建订单
	PARTIALLY_FILLED 部分成交
	FILLED 全部成交
	CANCELED 已撤销
	REJECTED 订单被拒绝
	EXPIRED 订单过期(根据timeInForce参数规则)


## 订单种类 (orderTypes, type):
	LIMIT 限价单
	MARKET 市价单
	STOP 止损限价单
	STOP_MARKET 止损市价单
	TAKE_PROFIT 止盈限价单
	TAKE_PROFIT_MARKET 止盈市价单
	TRAILING_STOP_MARKET 跟踪止损单

## 订单方向 (side):
	BUY 买入
	SELL 卖出

## 持仓方向:
	BOTH 单一持仓方向
	LONG 多头(双向持仓下)
	SHORT 空头(双向持仓下)

## 有效方式 (timeInForce):
	GTC - Good Till Cancel 成交为止
	IOC - Immediate or Cancel 无法立即成交(吃单)的部分就撤销
	FOK - Fill or Kill 无法全部立即成交就撤销
	GTX - Good Till Crossing 无法成为挂单方就撤销

## 条件价格触发类型 (workingType)
	MARK_PRICE
	CONTRACT_PRICE

## 响应类型 (newOrderRespType)
	ACK
	RESULT
*/

type crossType int

const (
	noCross      crossType = 0
	ascendCross  crossType = 1 // 上穿
	descendCross crossType = 2 // 下穿
)
