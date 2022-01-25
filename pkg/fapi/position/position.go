package position

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/account"
	"github.com/bitbeliever/binance-api/pkg/fapi/internal/principal"
	"github.com/bitbeliever/binance-api/pkg/fapi/order"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
	"math"
	"time"
)

// ClosePositionByOrderResp 平仓
/*
Open position:
	Long : positionSide=LONG, side=BUY
	Short: positionSide=SHORT, side=SELL

Close position:
	Close long position: positionSide=LONG, side=SELL
	Close short position: positionSide=SHORT, side=BUY
*/
func ClosePositionByOrderResp(o *futures.CreateOrderResponse) error {
	closeOrder, err := order.CreateOrderDual(o.Symbol, reverseSideType(o.Side), o.PositionSide, o.OrigQuantity)
	if err != nil {
		return err
	}

	log.Println("to close positions", helper.ToJson(closeOrder))
	return nil
}

func ClosePosition(position *futures.AccountPosition) error {
	var side futures.SideType
	var amt string
	if position.PositionAmt[0] == '-' {
		side = futures.SideTypeBuy
		amt = position.PositionAmt[1:]
	} else {
		side = futures.SideTypeSell
		amt = position.PositionAmt
	}
	_, err := order.CreateOrderDual(position.Symbol, side, position.PositionSide, amt)
	if err != nil {
		return err
	}

	return nil
}

/// todo
func openPosition() {}

// long/short reverse
func reversePositionSide(side futures.PositionSideType) futures.PositionSideType {
	if side == futures.PositionSideTypeLong {
		return futures.PositionSideTypeShort
	} else if side == futures.PositionSideTypeShort {
		return futures.PositionSideTypeLong
	} else {
		log.Println("wrong side")
		return futures.PositionSideTypeBoth
	}
}

// buy/sell reverse
func reverseSideType(sideType futures.SideType) futures.SideType {
	if sideType == futures.SideTypeBuy {
		return futures.SideTypeSell
	} else {
		return futures.SideTypeBuy
	}
}

func CloseAllPositions() {
	pos, err := account.QueryAccountPositions()
	if err != nil {
		panic(err)
	}

	for _, position := range pos {
		if err := ClosePosition(position); err != nil {
			log.Println(err)
		}
	}
}

func CloseAllPositionsBySymbol(symbol string) (err error) {
	pos, err := account.QueryAccountPositionsBySymbol(symbol)
	if err != nil {
		return err
	}

	for _, position := range pos {
		if err = ClosePosition(position); err != nil {
			log.Println(err)
		}
	}

	return err
}

// 仓位监控
func positionMonitor() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			pos, err := account.QueryAccountPositions()
			if err != nil {
				log.Println(err)
				continue
			}
			for _, p := range pos {
				pnl := helper.Str2Float64(p.UnrealizedProfit)
				if pnl < 0 && math.Abs(pnl) > principal.StopPNL() {
					log.Println("stop loss reach")
					if err := ClosePosition(p); err != nil {
						log.Println(err)
						return
					}
				}
			}
		}
	}
}
