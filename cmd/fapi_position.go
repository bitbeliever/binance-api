package main

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"log"
)

func main() {
	// sell position
	/*
		 created order {"symbol":"BNBUSDT","orderId":37122373487,"clientOrderId":"m3otKFUefMdusXF8dnikd3","price":"0","origQty":"0.10","executedQty":"0","cumQuote":"0","reduceOnly":false,"status":"
		NEW","stopPrice":"0","timeInForce":"GTC","type":"MARKET","side":"BUY","updateTime":1642168434842,"workingType":"CONTRACT_PRICE","activatePrice":"","priceRate":"","avgPrice":"0.00000","positionSide":"BOTH","closePosition":false,"priceProtect":false,"rateLimitOrder10s":"0","r
		ateLimitOrder1m":"1"}

	*/
	//fapi.CreateOrder(fapi.ETH, futures.SideTypeSell, "0.01")
	//fapi.QueryOpenOrders()

	fapi.QueryAccountBalance()
	fapi.QueryAccount()

	log.Println(fapi.CreateOrder(fapi.ETH, futures.SideTypeBuy, "0.001"))

	//log.Println(fapi.ModifyLeverage(fapi.ETH, 100))

	fapi.QueryAccountBalance()
	fapi.QueryAccount()
}
