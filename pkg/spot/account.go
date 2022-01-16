package spot

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"log"
	"time"
)

// AccountService /api/v3/account get account balances
func AccountService() (*binance.Account, error) {
	c := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	account, err := c.NewGetAccountService().Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//for _, balance := range account.Balances {
	//	if fapi.Str2Float64(balance.Locked) > 0 || fapi.Str2Float64(balance.Free) > 0 {
	//		log.Println(toJson(balance))
	//	}
	//}
	return account, nil
}

func AccountBalances() ([]binance.Balance, error) {
	account, err := AccountService()
	if err != nil {
		return nil, err
	}

	var balances []binance.Balance
	for _, b := range account.Balances {
		if str2Float64(b.Locked) != 0 || str2Float64(b.Free) != 0 {
			balances = append(balances, b)
		}
	}

	return balances, nil
}
