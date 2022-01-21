package fapi

import (
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"log"
)

// V1 穿过布林线
func v1Cross(symbol string, bRes bollResult, lastKline *futures.Kline) error {
	if isCrossingLine(bRes, lastKline) {
		// todo
		log.Println("crossing... current market", toJson(bRes), toJson(lastKline))
		//CreateOrder(symbol, futures.SideTypeBuy)
		return nil
	}

	return nil
}

// V2 区分上下穿==================
func distUpperLowerCross(symbol string, bRes bollResult, lastKline *futures.Kline) error {

	switch calCrossType(bRes, lastKline) {
	case ascendCross:
		log.Println("asc crossing... current market", toJson(bRes), toJson(lastKline))
		flog.Println("asc crossing... current market", toJson(bRes), toJson(lastKline))
		CreateOrder(symbol, futures.SideTypeBuy, "0.05") // qty todo
		//pinfo.Lock()
		//pinfo.position += 0.05
		//pinfo.Unlock()
		go monitor(symbol, lastKline.Close, 10, futures.SideTypeSell)
	case descendCross:
		log.Println("desc crossing... current market", toJson(bRes), toJson(lastKline))
		flog.Println("asc crossing... current market", toJson(bRes), toJson(lastKline))
		CreateOrder(symbol, futures.SideTypeSell, "0.05") // qty todo
		//pinfo.Lock()
		//pinfo.position -= 0.05
		//pinfo.Unlock()
		go monitor(symbol, lastKline.Close, 10, futures.SideTypeBuy)
	}
	return nil
}

type strategy interface {
	Do(symbol string, bRes bollResult, lastKline *futures.Kline) error
	StopLoss() error
	TakeProfit() error
}

type doubleOpenStrategy struct {
	opened bool

	upperChans []chan struct{}
	lowerChans []chan struct{}
}

// v3 中线双开策略, for now, 不能有未平的仓位
// V3 穿过中线, 双开 有未平的仓========================================
// todo yesterday: 达到中线后, 空/多判断, 开其中不存在的(可能多单或者空单只存其一)
func (s *doubleOpenStrategy) mbDoubleOpenPosition(symbol string, bRes bollResult, lastKline *futures.Kline) error {
	if !bollCross(bRes, lastKline) {
		return nil
	}

	//log.Println("crossing:", toJson(bRes), lastKline.Close)

	// 止损已经挂起, 直接返回
	if principal.stopLossFn != nil {
		return nil
	}

	//positions, err := QueryAccountPositions()
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}

	if bollCrossMB(bRes, lastKline) {
		if s.opened {
			return nil
		}
		// open position
		log.Println("达到中线======== to open position", toJson(bRes), lastKline.Close)
		// buy long
		orderBuy, err := CreateOrderDual(symbol, futures.SideTypeBuy, futures.PositionSideTypeLong, calcQty(principal.singleBetBalance(), lastKline.Close))
		if err != nil {
			log.Println(err)
			return err
		}
		// short sell
		orderSell, err := CreateOrderDual(symbol, futures.SideTypeSell, futures.PositionSideTypeShort, calcQty(principal.singleBetBalance(), lastKline.Close))
		if err != nil { // todo close
			log.Println(err)
			return err
		}

		log.Println("buy order", toJson(orderBuy))
		log.Println("sell order", toJson(orderSell))
		principal.orderLong = orderBuy
		principal.orderShort = orderSell
		s.opened = true
	} else if Str2Float64(lastKline.Close) >= bRes.UP { // 触碰上线 平多单 止盈

		// 止盈
		// 还未平仓
		if principal.orderLong != nil {
			if err := closePositionByOrderResp(principal.orderLong); err != nil {
				return err
			}
			principal.orderLong = nil
			log.Println("达到上线======== ", toJson(bRes), lastKline.Close)
		}
		// 止损 todo
		if principal.orderShort != nil {
			if err := closePositionByOrderResp(principal.orderShort); err != nil {
				return err
			}
			log.Println("达到下线======== ", toJson(bRes), lastKline.Close)
			principal.orderShort = nil
		}

		s.opened = false
	} else if Str2Float64(lastKline.Close) <= bRes.DN { // 触碰下线
		// 还未平仓
		if principal.orderShort != nil {
			if err := closePositionByOrderResp(principal.orderShort); err != nil {
				return err
			}
			principal.orderShort = nil
		}

		// 止损todo
		if principal.orderLong != nil {
			if err := closePositionByOrderResp(principal.orderLong); err != nil {
				return err
			}
			principal.orderLong = nil
		}

		s.opened = false
	}

	return nil
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
	ch := make(chan struct{})
	s.lowerChans = append(s.lowerChans, ch)
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

// by signal
func (s *doubleOpenStrategy) mbDoubleOpenPositionByChannel(symbol string, bRes bollResult, lastKline *futures.Kline) error {
	if !bollCross(bRes, lastKline) {
		return nil
	}

	if bollCrossMB(bRes, lastKline) {
		if s.opened {
			return nil
		}
		// open position
		log.Println("达到中线======== to open position", toJson(bRes), lastKline.Close)
		// buy long
		orderBuy, err := CreateOrderDual(symbol, futures.SideTypeBuy, futures.PositionSideTypeLong, calcQty(principal.singleBetBalance(), lastKline.Close))
		if err != nil {
			log.Println(err)
			return err
		}
		// short sell
		orderSell, err := CreateOrderDual(symbol, futures.SideTypeSell, futures.PositionSideTypeShort, calcQty(principal.singleBetBalance(), lastKline.Close))
		if err != nil { // todo close
			log.Println(err)
			return err
		}

		log.Println("buy order", toJson(orderBuy))
		log.Println("sell order", toJson(orderSell))
		go monitorOrderTP(orderBuy, s.subscribeUpper(), s.subscribeLower())  // 设置止盈
		go monitorOrderTP(orderSell, s.subscribeUpper(), s.subscribeLower()) // 设置止盈

		pos, err := QueryAccountPositions()
		if err != nil {
			log.Println(err)
			return err
		}
		for _, p := range pos {
			if p.PositionSide == futures.PositionSideTypeLong { // 多单, stop by upper
				go watchPNLStopLimit(p, s.subscribeUpper()) // 设置止损
			} else if p.PositionSide == futures.PositionSideTypeShort { // 空单, stop by lower
				go watchPNLStopLimit(p, s.subscribeLower()) // 设置止损
			}
		}

		s.opened = true
	} else if Str2Float64(lastKline.Close) >= bRes.UP-0.5 { // 触碰上线 平多单 止盈
		s.pubUpper()
		s.opened = false
	} else if Str2Float64(lastKline.Close) <= bRes.DN+0.5 { // 触碰下线
		s.pubLower()
		s.opened = false
	}

	return nil
}

func Test() {
	log.Println(calcQty(10, "3136.32"))
	log.Println(10 / 3136.32)
}
