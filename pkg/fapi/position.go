package fapi

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"math"
	"time"
)

// 平仓
/*
Open position:
	Long : positionSide=LONG, side=BUY
	Short: positionSide=SHORT, side=SELL

Close position:
	Close long position: positionSide=LONG, side=SELL
	Close short position: positionSide=SHORT, side=BUY
*/
func closePositionByOrderResp(order *futures.CreateOrderResponse) error {
	//closeOrder, err := CreateOrderDual(order.Symbol, futures.SideTypeSell, reversePositionSide(order.PositionSide), order.OrigQuantity)
	//closeOrder, err := CreateOrderDual(order.Symbol, futures.SideTypeSell, order.PositionSide, order.OrigQuantity)
	closeOrder, err := CreateOrderDual(order.Symbol, reverseSideType(order.Side), order.PositionSide, order.OrigQuantity)
	if err != nil {
		return err
	}
	//if positionSide == futures.PositionSideTypeLong {
	//	//principal.closeLongOrder = closeOrder
	//} else if positionSide == futures.PositionSideTypeShort {
	//	//principal.closeShortOrder = closeOrder
	//} else {
	//	return fmt.Errorf("wrong side %s", positionSide)
	//}

	log.Println("to close positions", toJson(closeOrder))
	return nil
}

func closePosition(position *futures.AccountPosition) error {
	var side futures.SideType
	var amt string
	if position.PositionAmt[0] == '-' {
		side = futures.SideTypeBuy
		amt = position.PositionAmt[1:]
	} else {
		side = futures.SideTypeSell
		amt = position.PositionAmt
	}
	_, err := CreateOrderDual(position.Symbol, side, position.PositionSide, amt)
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

// 用户持仓风险V2 /fapi/v2/positionRisk
func positionRisk(symbol string) ([]*futures.PositionRisk, error) {
	return NewClient().NewGetPositionRiskService().Symbol(symbol).Do(context.Background())
}

func CloseAllPositions() {
	pos, err := QueryAccountPositions()
	if err != nil {
		panic(err)
	}

	for _, position := range pos {
		if err := closePosition(position); err != nil {
			log.Println(err)
		}
	}
}

// 仓位监控
func positionMonitor() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			pos, err := QueryAccountPositions()
			if err != nil {
				log.Println(err)
				continue
			}
			for _, p := range pos {
				pnl := Str2Float64(p.UnrealizedProfit)
				if pnl < 0 && math.Abs(pnl) > principal.stopPNL() {
					log.Println("stop loss reach")
					if err := closePosition(p); err != nil {
						log.Println(err)
						return
					}
				}
			}
		}
	}
}
