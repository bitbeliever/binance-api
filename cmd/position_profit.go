package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi/position"
	"log"
	"time"
)

func main() {
	tick := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-tick.C:
			profitSum, err := position.UnrealizedProfit()
			if err != nil {
				log.Println(err)
				continue
			}

			if profitSum != 0 {
				log.Println("un profitSum:", profitSum)
			}
		}
	}
}
