package main

import (
	"github.com/bitbeliever/binance-api/pkg/account"
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/fapi/order"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
	"strconv"
)

func main() {
	//var keys []string
	//for i := 1; i <= 10; i++ {
	//	keys = append(keys, keyPhase(i), keyPhase(-i))
	//}
	//if err := cache.Client.Del(keys...).Err(); err != nil {
	//	log.Println(err)
	//}
	//return

	helper.JsonLog(account.QueryAccountPositions())
	//fapi.QueryAccount()
	//fapi.QueryAllOrders(fapi.BNB)
	o, err := order.DualBuyLong(fapi.BCH, "0.5")
	if err != nil {
		log.Println(err)
	}
	log.Println(helper.ToJson(o))

}

func keyPhase(i int) string {
	return "smooth_" + strconv.Itoa(i)
}
