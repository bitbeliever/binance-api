package strategy

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/fapi/internal/indicator"
	"github.com/bitbeliever/binance-api/pkg/fapi/internal/principal"
	"github.com/bitbeliever/binance-api/pkg/fapi/order"
	"github.com/bitbeliever/binance-api/pkg/fapi/position"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

type pyramid struct {
	byArithmeticProgression bool
	byGeometricProgression  bool
	segments                int // 段数

	//gap                     float64
	//firstHalfGap  float64
	//secondHalfGap float64

	layers int // 金字塔层数
}

//func (p pyramid) phase2(price string, MB float64) int {
//	curPrice := helper.Str2Float64(price)
//
//	for i := 1; i <= p.segments; i++ {
//		if float64(i)*p.gap+MB > curPrice {
//			return i
//		}
//	}
//	return 0
//}

func (p pyramid) calcGap(b indicator.Boll) (float64, float64) {
	res := b.Result()

	return (res.UP - res.MB) / float64(p.segments), (res.MB - res.DN) / float64(p.segments)
}

func (p pyramid) phase(b indicator.Boll) int {
	pri := helper.Str2Float64(b.CurrentPrice())
	// 大于上线, 低于下线
	if pri > b.Result().UP || pri < b.Result().DN {
		return 0
	}

	firstHalfGap, secondHalfGap := p.calcGap(b)
	if b.IsFirstHalf() {
		for i := 1; i <= p.segments; i++ {
			if float64(i)*firstHalfGap+b.Result().MB > pri {
				return i
			}
		}
	} else if b.IsSecondHalf() {
		for i := 1; i <= p.segments; i++ {
			if b.Result().MB-float64(i)*secondHalfGap < pri {
				return -i
			}
		}
	} else {
		panic("ERR wrong boll" + helper.ToJson(b.Klines()))
	}

	panic("ERR gap boll" + helper.ToJson(b.Klines()))
}

func newPyramid(byArithmetic bool) pyramid {

	p := pyramid{
		byArithmeticProgression: byArithmetic,  // 算术 等差
		byGeometricProgression:  !byArithmetic, // 几何 等比
		segments:                20,            // 段数设为20
	}
	//bRes := boll.Result()
	//p.firstHalfGap = (bRes.UP - bRes.MB) / float64(p.segments)
	//p.secondHalfGap = (bRes.MB - bRes.DN) / float64(p.segments)

	//if isFirstHalf {
	//	p.gap = (bRes.UP - bRes.MB) / float64(p.segments)
	//} else {
	//	p.gap = (bRes.MB - bRes.DN) / float64(p.segments)
	//}

	return p
}

type Smooth struct {
	symbol string
	p      pyramid

	//firstHalfPyramid  pyramid
	//secondHalfPyramid pyramid
	opened bool

	initLongOrder  *futures.CreateOrderResponse
	initShortOrder *futures.CreateOrderResponse
	upperCh        chan struct{}
	lowerCh        chan struct{}

	//addedLongAmt  float64
	//addedShortAmt float64
	leverage *futures.SymbolLeverage

	state map[int]float64
}

func NewSmooth(symbol string, boll indicator.Boll) *Smooth {
	s := &Smooth{
		symbol: symbol,
		p:      newPyramid(true),
		//firstHalfPyramid:  newPyramid(true),
		//secondHalfPyramid: newPyramid(true),
		state: make(map[int]float64),
	}

	//go s.monitorUPDN()
	return s
}

func (s *Smooth) Do(symbol string, boll indicator.Boll) error {
	// 跨中线 双开
	if boll.CrossMB() {
		// 平掉开了的仓位
		s.reset()
		if s.opened { // todo

		}

		longOrder, err := order.DualBuyLong(symbol, calcQty(principal.SingleBetBalance(), boll.LastKline().Close, s.leverage.Leverage))
		if err != nil {
			return err
		}
		log.Println("中线 long order", helper.ToJson(longOrder))
		shortOrder, err := order.DualSellShort(symbol, calcQty(principal.SingleBetBalance(), boll.LastKline().Close, s.leverage.Leverage))
		if err != nil {
			return err
		}
		log.Println("中线 short order", helper.ToJson(shortOrder))

		s.opened = true
		s.initLongOrder = longOrder
		s.initShortOrder = shortOrder
	} else if boll.IsFirstHalf() { // 上半段
		phase := s.p.phase(boll)
		if _, ok := s.state[phase]; !ok {
			o, err := order.DualSellShort(s.symbol, principal.Qty())
			if err != nil {
				return err
			}
			log.Println("上半段开仓", helper.ToJson(o))
			s.state[phase] = 1
			log.Println(helper.ToJson(s.state))
		}

	} else if boll.IsSecondHalf() { // 下半段
		phase := s.p.phase(boll)
		if _, ok := s.state[phase]; !ok {
			o, err := order.DualBuyLong(s.symbol, principal.Qty())
			if err != nil {
				return err
			}
			log.Println("下半段开仓", helper.ToJson(o))
			s.state[phase] = 1
			log.Println(helper.ToJson(s.state))
		}
	} else if boll.CrossUP() { // 触碰上线
		if s.initLongOrder != nil {
			err := position.ClosePositionByOrderResp(s.initLongOrder)
			if err != nil {
				return nil
			}
		}
	} else if boll.CrossDN() { // 触碰下线
		if s.initShortOrder != nil {
			err := position.ClosePositionByOrderResp(s.initShortOrder)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Smooth) CloseBaseLongShort() {
	// 平所以仓
	if err := position.CloseAllPositionsBySymbol(s.symbol); err != nil {
		log.Println(err)
		return
	}

}

func (s *Smooth) reset() {
	if err := position.CloseAllPositionsBySymbol(s.symbol); err != nil {
		log.Println(err)
	}

	s.state = make(map[int]float64)
	s.opened = false
}

func (s *Smooth) monitorUPDN() {
	for {
		select {
		case <-s.upperCh:
			if s.initLongOrder != nil {
				err := position.ClosePositionByOrderResp(s.initLongOrder)
				if err != nil {
					log.Println(err)
				}
				return
			}
		case <-s.lowerCh:
			if s.initShortOrder != nil {
				err := position.ClosePositionByOrderResp(s.initShortOrder)
				if err != nil {
					log.Println(err)
				}
				return
			}
		}
	}
}
