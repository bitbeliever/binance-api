package fapi

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"log"
)

// 平仓
func closePosition(order *futures.CreateOrderResponse, positionSide futures.PositionSideType) error {
	//order.
	//closeOrder, err := CreateOrderBothSide(order.Symbol, futures.SideTypeSell, reversePositionSide(order.PositionSide), order.OrigQuantity)
	closeOrder, err := CreateOrderBothSide(order.Symbol, futures.SideTypeSell, order.PositionSide, order.OrigQuantity)
	if err != nil {
		return err
	}
	if positionSide == futures.PositionSideTypeLong {
		principal.closeLong = closeOrder
	} else if positionSide == futures.PositionSideTypeShort {
		principal.closeShort = closeOrder
	} else {
		return fmt.Errorf("wrong side %s", positionSide)
	}

	log.Println("to close positions", toJson(closeOrder))
	return nil
}

/// todo
func openPosition() {}

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

// 用户持仓风险V2 /fapi/v2/positionRisk
func positionRisk(symbol string) ([]*futures.PositionRisk, error) {
	return NewClient().NewGetPositionRiskService().Symbol(symbol).Do(context.Background())
	//risks, err :=
	//if err != nil {
	//	return
	//}
}
