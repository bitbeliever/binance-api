package fur

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"time"
)

func orderTest() {
	client := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	order, err := client.NewCreateOrderService().Symbol("BNBETH").
		Side(futures.SideTypeBuy).Type(futures.OrderTypeLimit).
		TimeInForce(futures.TimeInForceTypeGTC).Quantity("5").
		Price("0.0030000").Do(ctx)
	if err != nil {
		log.Println(err)
	}

	log.Println(toJson(order))
}
