package spot

import (
	"context"
	"github.com/bitbeliever/binance-api/pkg/fur"
	"log"
	"time"
)

// AccountServiceTest get account balances
func AccountServiceTest() {
	c := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	res, err := c.NewGetAccountService().Do(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	//log.Println(res)
	for _, balance := range res.Balances {
		if fur.Str2Float64(balance.Locked) > 0 || fur.Str2Float64(balance.Free) > 0 {
			log.Println(toJson(balance))
		}
	}
}
