package spot

import (
	"context"
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"log"
	"time"
)

// AccountService /api/v3/account get account balances
func AccountService() {
	c := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.NewGetAccountService().Do(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	//log.Println(res)
	for _, balance := range res.Balances {
		if fapi.Str2Float64(balance.Locked) > 0 || fapi.Str2Float64(balance.Free) > 0 {
			log.Println(toJson(balance))
		}
	}
}

func QueryBalance() {
}
