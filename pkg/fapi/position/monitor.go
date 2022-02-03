package position

import (
	"github.com/bitbeliever/binance-api/pkg/account"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
	"time"
)

func MonitorPositions(symbol string, tl float64, tick time.Duration) chan float64 {
	ch := make(chan float64, 2<<10)
	go func() {
		tick := time.NewTicker(tick)
		for {
			select {
			case <-tick.C:
				positions, err := account.QueryAccountPositionsBySymbol(symbol)
				if err != nil {
					log.Println(err)
					continue
				}

				var profitSum float64
				for _, p := range positions {
					profitSum += helper.Str2Float64(p.UnrealizedProfit)
				}

				if profitSum <= tl {
					profit, err := CloseAllPositionsBySymbol(symbol)
					if err != nil {
						log.Println(err)
						continue
					}
					log.Println("monitor close profit", profit)
				}
			}

		}
	}()
	return ch
}
