package fapi

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"time"
)

// QueryAccount 账户查询 账户信息 fapi/v1/account
func QueryAccount() *futures.Account {
	c := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	res, err := c.NewGetAccountService().Do(ctx)
	if err != nil {
		log.Println(err)
		return res
	}

	//log.Println(toJson(res))
	for _, asset := range res.Assets {
		if Str2Float64(asset.WalletBalance) > 0 {
			log.Println(toJson(asset))
		}
	}
	for _, position := range res.Positions {
		if Str2Float64(position.PositionAmt) > 0 {
			log.Println(toJson(position))
		}

	}
	//for _, balance := range res.Balances {
	//	if Str2Float64(balance.Locked) > 0 || Str2Float64(balance.Free) > 0 {
	//		log.Println(balance)
	//	}
	//}

	return res
}

// QueryAccountBalance 账户余额
func QueryAccountBalance() {
	balances, err := NewClient().NewGetBalanceService().Do(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(toJsonIndent(balances))
}
