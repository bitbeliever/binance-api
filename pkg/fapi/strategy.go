package fapi

import (
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"log"
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

	longOrder  *futures.CreateOrderResponse
	shortOrder *futures.CreateOrderResponse
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

			go s.monitorOrderSL(longOrder)
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

			go s.monitorOrderSL(shortOrder)
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
