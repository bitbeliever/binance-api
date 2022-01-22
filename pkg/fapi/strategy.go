package fapi

import (
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"math"
)

type strategy interface {
	Do(symbol string, bRes bollResult, lastKline *futures.Kline) error
	StopLoss() error
	TakeProfit() error
}

type doubleOpenStrategy struct {
	upperChans []chan struct{}
	lowerChans []chan struct{}
	stopChans  []chan struct{}

	// todo add lock
	longOrder        *futures.CreateOrderResponse
	longOrderStopCh  chan struct{}
	shortOrder       *futures.CreateOrderResponse
	shortOrderStopCh chan struct{}
}

func newDoubleOpenStrategy() *doubleOpenStrategy {
	return &doubleOpenStrategy{
		longOrderStopCh:  make(chan struct{}, 256),
		shortOrderStopCh: make(chan struct{}, 256),
	}
}

func calcQty(spend float64, closeStr string) string {
	price := Str2Float64(closeStr)
	return fmt.Sprintf("%.2f", spend/price)
	//spend / price
	//return strconv.FormatFloat(math.Round(spend/price*100)/100, 'f', 10, 64)
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

// by signal
func (s *doubleOpenStrategy) mbDoubleOpenPositionByChannel(symbol string, bRes bollResult, lastKline *futures.Kline) error {
	if !bollCross(bRes, lastKline) {
		return nil
	}

	if bollCrossMB(bRes, lastKline) {

		// buy long
		if s.longOrder == nil {

			longOrder, err := CreateOrderDual(symbol, futures.SideTypeBuy, futures.PositionSideTypeLong, calcQty(principal.singleBetBalance(), lastKline.Close))
			if err != nil {
				log.Println(err)
				return err
			}
			s.longOrder = longOrder
			log.Println("达到中线=== open long order", toJson(bRes), lastKline.Close, toJson(longOrder))

			go s.monitorOrderSL(longOrder, s.longOrderStopCh)
		}
		if s.shortOrder == nil {
			// short sell
			shortOrder, err := CreateOrderDual(symbol, futures.SideTypeSell, futures.PositionSideTypeShort, calcQty(principal.singleBetBalance(), lastKline.Close))
			if err != nil { // todo close
				log.Println(err)
				return err
			}
			s.shortOrder = shortOrder
			log.Println("达到中线=== open short order", toJson(bRes), lastKline.Close, toJson(shortOrder))

			go s.monitorOrderSL(shortOrder, s.shortOrderStopCh)
		}

	} else if Str2Float64(lastKline.Close) >= bRes.UP-0.5 { // 触碰上线 平多单 止盈
		s.pubUpper()
	} else if Str2Float64(lastKline.Close) <= bRes.DN+0.5 { // 触碰下线
		s.pubLower()
	}

	return nil
}

func Test() {
	log.Println(calcQty(10, "3136.32"))
	log.Println(10 / 3136.32)
}

func (s *doubleOpenStrategy) monitorOrderTP(chUpper, chLower chan struct{}) {
	for {
		select {
		case <-chUpper:
			if s.longOrder != nil {
				log.Println("触发多单止盈")
				if err := closePositionByOrderResp(s.longOrder); err != nil {
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
				if err := closePositionByOrderResp(s.shortOrder); err != nil {
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
	ch, err := AggTradePrice(p.Symbol)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		select {
		case curPriceStr := <-ch:
			pnl := calcPNL(Str2Float64(p.OrigQuantity), p.PositionSide, Str2Float64(p.AvgPrice), Str2Float64(curPriceStr))
			//log.Println("current pnl", pnl)
			// 触发止损
			if pnl < 0 && math.Abs(pnl) > principal.stopPNL() {
				log.Printf("触发止损 pnl: %v stopPNL: %v, pSide %v \n", pnl, principal.stopPNL(), p.PositionSide)
				if err := closePositionByOrderResp(p); err != nil {
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
