package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/helper"
)

func main() {
	// sell position
	/*
			 created order {"symbol":"BNBUSDT","orderId":37122373487,"clientOrderId":"m3otKFUefMdusXF8dnikd3","price":"0","origQty":"0.10","executedQty":"0","cumQuote":"0","reduceOnly":false,"status":"
			NEW","stopPrice":"0","timeInForce":"GTC","type":"MARKET","side":"BUY","updateTime":1642168434842,"workingType":"CONTRACT_PRICE","activatePrice":"","priceRate":"","avgPrice":"0.00000","positionSide":"BOTH","closePosition":false,"priceProtect":false,"rateLimitOrder10s":"0","r
			ateLimitOrder1m":"1"}

		[{"isolated":false,"leverage":"20","initialMargin":"2.47830000","maintMargin":"0.32217900","openOrderInitialMargin":"0","positionInitialMargin":"2.47830000","symbol":"BNBUSDT","unrealiz
		edProfit":"0.00000000","entryPrice":"495.66","maxNotional":"250000","positionSide":"BOTH","positionAmt":"0.10","notional":"49.56600000","isolatedWallet":"0","updateTime":1642327506307}]

	*/

	//helper.JsonLog(fapi.QueryBalance())
	//helper.JsonLog(fapi.QueryAccountAssets())
	//helper.JsonLog(fapi.QueryAccountPositions())

	//helper.JsonLog(fapi.ModifyLeverage(fapi.BNB, 75))
	//helper.JsonLog(fapi.CreateOrder(fapi.BNB, futures.SideTypeSell, "0.05"))
	//helper.JsonLog(fapi.ModifyLeverage(fapi.ETH, 100))

	helper.JsonLog(fapi.QueryBalance())
	helper.JsonLog(fapi.QueryAccountAssets())
	helper.JsonLog(fapi.QueryAccountPositions())
}
