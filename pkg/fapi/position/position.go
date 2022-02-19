package position

import (
	"fmt"
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

func CloseOrderResps(orders []*futures.CreateOrderResponse) (err error) {
	if len(orders) == 0 {
		return nil
	}
	var amt float64
	for _, o := range orders {
		amt += helper.Str2Float64(o.OrigQuantity)
	}
	orders[0].OrigQuantity = helper.FloatToStr(amt)
	err = ClosePositionByOrderResp(orders[0])
	//profit = order.CalcPNL(orders[0])
	return
}

// ClosePositionBySymbol todo test
func ClosePositionBySymbol(symbol string, amt float64) error {
	var side futures.SideType
	var amtStr = helper.FloatToStr(amt)
	var pSide futures.PositionSideType
	if amt > 0 {
		side = futures.SideTypeSell
		pSide = futures.PositionSideTypeLong
	} else {
		side = futures.SideTypeBuy
		pSide = futures.PositionSideTypeShort
	}

	o, err := order.CreateOrderDual(symbol, side, pSide, amtStr)
	_ = o
	if err != nil {
		return err
	}
	return err
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

func CloseAllPositions() (profit float64, err error) {
	pos, err := account.QueryAccountPositions()
	if err != nil {
		return
	}

	for _, position := range pos {
		profit += helper.Str2Float64(position.UnrealizedProfit)
		if er := ClosePosition(position); err != nil {
			err = fmt.Errorf("%w", er)
		}
	}
	return
}

func CloseAllPositionsBySymbol(symbol string) (float64, error) {
	pos, err := account.QueryAccountPositionsBySymbol(symbol)
	if err != nil {
		return 0, err
	}

	var sum float64
	for _, position := range pos {
		sum += helper.Str2Float64(position.UnrealizedProfit)
		if err = ClosePosition(position); err != nil {
			log.Println(err)
		}
	}

	return sum, nil
}

func ClosePositionBySymbolSide(symbol string, pSide futures.PositionSideType) (float64, error) {
	pos, err := account.QueryAccountPositionsBySymbol(symbol)
	if err != nil {
		return 0, err
	}

	var sum float64
	for _, position := range pos {
		if position.PositionSide == pSide {
			if err = ClosePosition(position); err != nil {
				log.Println(err)
			}
			sum += helper.Str2Float64(position.UnrealizedProfit)
		}
	}

	return sum, nil
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

type profit struct {
	Symbol       string
	PositionSide string
	Amt          string
	Profit       string
}

func (p profit) String() string {
	return fmt.Sprintf("%s %s %s %s\n", p.Symbol, p.Amt, p.PositionSide, p.Profit)
}

func UnrealizedProfit() (float64, string, []*futures.AccountPosition, error) {
	pos, err := account.QueryAccountPositions()
	if err != nil {
		return 0, "", nil, err
	}
	var sum float64
	var detail = "\n"
	var profits []profit
	for _, p := range pos {
		sum += helper.Str2Float64(p.UnrealizedProfit)
		detail += fmt.Sprintf("%s %s %s %s\n", p.Symbol, p.PositionAmt, p.PositionSide, p.UnrealizedProfit)
		profits = append(profits, profit{
			Symbol:       p.Symbol,
			PositionSide: string(p.PositionSide),
			Amt:          p.PositionAmt,
			Profit:       p.UnrealizedProfit,
		})
	}
	detail += fmt.Sprintf("total profit %v\n", sum)
	return sum, detail, pos, nil
}

func UnrealizedProfitSymbol(symbol string) (float64, error) {
	pos, err := account.QueryAccountPositionsBySymbol(symbol)
	if err != nil {
		return 0, err
	}
	var sum float64
	for _, p := range pos {
		sum += helper.Str2Float64(p.UnrealizedProfit)
	}
	return sum, nil
}

func positionHandler(o *futures.CreateOrderResponse) {

}
