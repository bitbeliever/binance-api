package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi/position"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
	"time"
)

func main() {
	m := NewMotion()
	m.Reg("LTCUSDTLONG", tuple{1, -2})
	m.Run()
}

type tuple struct {
	tp float64
	sl float64
}

type motion struct {
	reg map[string]tuple
}

func (m *motion) Reg(s string, t tuple) {
	m.reg[s] = t
}

func NewMotion() *motion {
	return &motion{
		reg: make(map[string]tuple),
	}
}

func (m *motion) Run() {
	for k, v := range m.reg {
		log.Println(k, v)
	}
	tick := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-tick.C:
			profitSum, detail, profits, err := position.UnrealizedProfit()
			if err != nil {
				log.Println(err)
				continue
			}

			if profitSum != 0 {
				log.Println(detail)
				//log.Println("un profitSum:", profitSum)
			}

			for _, p := range profits {
				t, ok := m.reg[p.Symbol+string(p.PositionSide)]
				if !ok {
					continue
				}
				if helper.Str2Float64(p.UnrealizedProfit) > t.tp {
					log.Println("tp reached", t.tp, p.UnrealizedProfit)
					if err := position.ClosePosition(p); err != nil {
						log.Println(err)
					}
				} else if helper.Str2Float64(p.UnrealizedProfit) < t.sl {
					log.Println("sl reached", t.sl, p.UnrealizedProfit, p.EntryPrice)
					if err := position.ClosePosition(p); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}
