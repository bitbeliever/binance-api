package fur

import (
	"context"
	"log"
	"time"
)

func AccountServiceTest() {
	c := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	res, err := c.NewGetAccountService().Do(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(toJson(res))
	//for _, balance := range res.Balances {
	//	if Str2Float64(balance.Locked) > 0 || Str2Float64(balance.Free) > 0 {
	//		log.Println(balance)
	//	}
	//}
}
