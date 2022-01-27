package strategy

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/fapi/internal/indicator"
	"github.com/bitbeliever/binance-api/pkg/fapi/internal/principal"
	"github.com/bitbeliever/binance-api/pkg/fapi/order"
	"github.com/bitbeliever/binance-api/pkg/fapi/position"
	"github.com/bitbeliever/binance-api/pkg/fapi/trade"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
	"math"
)

type strategy interface {
	Do(symbol string, boll indicator.Boll) error
	StopLoss() error
	TakeProfit() error
}

// todo add lock
type doubleOpenStrategy struct {
	symbol     string
	upperChans []chan struct{}
	lowerChans []chan struct{}
	stopChans  []chan struct{}

	longOrder        *futures.CreateOrderResponse
	longOrderStopCh  chan struct{}
	shortOrder       *futures.CreateOrderResponse
	shortOrderStopCh chan struct{}
	leverage         *futures.SymbolLeverage
}

func NewDoubleOpenStrategy(symbol string, lev *futures.SymbolLeverage) *doubleOpenStrategy {
	s := &doubleOpenStrategy{
		symbol:           symbol,
		longOrderStopCh:  make(chan struct{}, 256),
		shortOrderStopCh: make(chan struct{}, 256),
		leverage:         lev,
	}

	// 设置止盈 call only once
	go s.monitorOrderTP(s.subscribeUpper(), s.subscribeLower())

	return s
}

func (s *doubleOpenStrategy) subscribeUpper() chan struct{} {
	ch := make(chan struct{}, 256)
	s.upperChans = append(s.upperChans, ch)
	return ch
}

func (s *doubleOpenStrategy) subscribeLower() chan struct{} {
	ch := make(chan struct{}, 256)
	s.lowerChans = append(s.lowerChans, ch)
	return ch
}

func (s *doubleOpenStrategy) subscribeStop() chan struct{} {
	ch := make(chan struct{}, 256)
	s.stopChans = append(s.stopChans, ch)
	return ch
}

func (s *doubleOpenStrategy) pubUpper() {
	for _, ch := range s.upperChans {
		ch <- struct{}{}
	}
}

func (s *doubleOpenStrategy) pubLower() {
	for _, ch := range s.lowerChans {
		ch <- struct{}{}
	}
}
func (s *doubleOpenStrategy) pubStop() {
	for _, ch := range s.stopChans {
		ch <- struct{}{}
	}
}

func (s *doubleOpenStrategy) Do(symbol string, boll indicator.Boll) error {
	if boll.Cross() {
		return nil
	}

	if boll.CrossMB() {
		// buy long
		if s.longOrder == nil {

			//longOrder, err := order.CreateOrderDual(symbol, futures.SideTypeBuy, futures.PositionSideTypeLong, calcQty(principal.SingleBetBalance(), boll.LastKline().Close, s.leverage.Leverage))
			longOrder, err := order.CreateOrderDual(symbol, futures.SideTypeBuy, futures.PositionSideTypeLong, principal.Qty())
			if err != nil {
				log.Println(err)
				return err
			}
			s.longOrder = longOrder
			log.Println("达到中线=== open long order", helper.ToJson(boll.Result()), boll.LastKline().Close, helper.ToJson(longOrder))

			go s.monitorOrderSL(longOrder, s.longOrderStopCh)
		}
		if s.shortOrder == nil {
			// short sell
			//shortOrder, err := order.CreateOrderDual(symbol, futures.SideTypeSell, futures.PositionSideTypeShort, calcQty(principal.SingleBetBalance(), boll.LastKline().Close, s.leverage.Leverage))
			shortOrder, err := order.CreateOrderDual(symbol, futures.SideTypeSell, futures.PositionSideTypeShort, principal.Qty())
			if err != nil { // todo close
				log.Println(err)
				return err
			}
			s.shortOrder = shortOrder
			log.Println("达到中线=== open short order", helper.ToJson(boll.Result()), boll.LastKline().Close, helper.ToJson(shortOrder))

			go s.monitorOrderSL(shortOrder, s.shortOrderStopCh)
		}

	} else if boll.CrossUP() { // 触碰上线 平多单 止盈
		s.pubUpper()
	} else if boll.CrossDN() { // 触碰下线
		s.pubLower()
	}

	return nil

}

func Test() {
	//log.Println(calcQty(10, "3136.32"))
	log.Println(10 / 3136.32)
}

func (s *doubleOpenStrategy) monitorOrderTP(chUpper, chLower chan struct{}) {
	for {
		select {
		case <-chUpper:
			if s.longOrder != nil {
				log.Println("触发多单止盈")
				if err := position.ClosePositionByOrderResp(s.longOrder); err != nil {
					log.Println(err)
					return
				}
				s.longOrder = nil
				s.longOrderStopCh <- struct{}{}
				return
			}

		case <-chLower:
			if s.shortOrder != nil {
				log.Println("触发空单止盈")
				if err := position.ClosePositionByOrderResp(s.shortOrder); err != nil {
					log.Println(err)
					return
				}
				s.shortOrder = nil
				s.shortOrderStopCh <- struct{}{}
				return
			}
		}

	}
}

func (s *doubleOpenStrategy) monitorOrderSL(p *futures.CreateOrderResponse, stop chan struct{}) {
	ch, err := trade.AggTradePrice(p.Symbol)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		select {
		case curPriceStr := <-ch:
			pnl := calcPNL(helper.Str2Float64(p.OrigQuantity), p.PositionSide, helper.Str2Float64(p.AvgPrice), helper.Str2Float64(curPriceStr))
			//log.Println("current pnl", pnl)
			// 触发止损
			if pnl < 0 && math.Abs(pnl) > principal.StopPNL() {
				log.Printf("触发止损 pnl: %v stopPNL: %v, pSide %v \n", pnl, principal.StopPNL(), p.PositionSide)
				if err := position.ClosePositionByOrderResp(p); err != nil {
					log.Println(err)
					return
				}

				// 消除monitor take profit
				if p.PositionSide == futures.PositionSideTypeLong && s.longOrder != nil {
					s.longOrder = nil
				} else if p.PositionSide == futures.PositionSideTypeShort && s.shortOrder != nil {
					s.shortOrder = nil
				}
				return
			}
		case <-stop:
			return
		}
	}
}
