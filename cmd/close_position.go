package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
)

func main() {
	// sell position
	/*
		 created order {"symbol":"BNBUSDT","orderId":37122373487,"clientOrderId":"m3otKFUefMdusXF8dnikd3","price":"0","origQty":"0.10","executedQty":"0","cumQuote":"0","reduceOnly":false,"status":"
		NEW","stopPrice":"0","timeInForce":"GTC","type":"MARKET","side":"BUY","updateTime":1642168434842,"workingType":"CONTRACT_PRICE","activatePrice":"","priceRate":"","avgPrice":"0.00000","positionSide":"BOTH","closePosition":false,"priceProtect":false,"rateLimitOrder10s":"0","r
		ateLimitOrder1m":"1"}

	*/
	//fapi.CreateOrder(fapi.BNB, futures.SideTypeSell, "0.3")
	//fapi.QueryOpenOrders()

	fapi.QueryAccountBalance()
	fapi.QueryAccount()
	fapi.QueryOpenOrders()
}
