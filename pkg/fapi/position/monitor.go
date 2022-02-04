package position

import (
	"log"
	"time"
)

func MonitorPositions(symbol string, stopLossLine float64, tick time.Duration, fn func()) {
	//ch := make(chan float64, 2<<10)
	go func() {
		tick := time.NewTicker(tick)
		for {
			select {
			case <-tick.C:
				profitSum, err := UnrealizedProfitSymbol(symbol)
				if err != nil {
					log.Println(err)
					continue
				}

				if profitSum <= stopLossLine {
					log.Printf("monitor sl profitSum: %v, SL: %v \n", profitSum, stopLossLine)
					fn()
				}
			}

		}
	}()
	//return ch
}
