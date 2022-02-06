package strategy

import (
	"encoding/json"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/configs"
	"github.com/bitbeliever/binance-api/pkg/cache"
	"github.com/bitbeliever/binance-api/pkg/fapi/indicator"
	"github.com/bitbeliever/binance-api/pkg/fapi/internal/principal"
	"github.com/bitbeliever/binance-api/pkg/fapi/order"
	"github.com/bitbeliever/binance-api/pkg/fapi/position"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
	"strconv"
	"time"
)

const (
	prefix  = "smooth_"
	pattern = prefix + "*"
)

type pyramid struct {
	byArithmeticProgression bool
	byGeometricProgression  bool
	segments                int // 段数

	//gap                     float64
	//firstHalfGap  float64
	//secondHalfGap float64
	//layers int // 金字塔层数
}

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
		byArithmeticProgression: byArithmetic,                         // 算术 等差
		byGeometricProgression:  !byArithmetic,                        // 几何 等比
		segments:                configs.Cfg.Strategy.PyramidSegments, // 段数设为20
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
	symbol         string
	p              pyramid
	opened         bool
	initLongOrder  *futures.CreateOrderResponse
	initShortOrder *futures.CreateOrderResponse

	//state    map[int]float64
	phaseOrders []*futures.CreateOrderResponse
}

func NewSmooth(symbol string) *Smooth {
	keys := cache.Client.Keys("smooth_*").Val()
	log.Println("init redis keys:", keys)

	o1, o2, err := resumeInitOrder()
	if err != nil {
		log.Println(err)
	}

	s := &Smooth{
		symbol:         symbol,
		p:              newPyramid(true),
		initLongOrder:  o1,
		initShortOrder: o2,
	}
	if len(keys) > 0 {
		s.opened = true
	}
	if o1 != nil || o2 != nil {
		s.opened = true
	}

	//go position.MonitorPositions(symbol, configs.Cfg.StopLoss, time.Second*2, func() {
	//	profit, err := position.CloseAllPositionsBySymbol(symbol)
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//	log.Println("monitor close profit", profit)
	//})
	go position.MonitorPositions(symbol, configs.Cfg.StopLoss, time.Second*2, s.reset)

	return s
}

func (s *Smooth) Do(lines []*futures.Kline) error {
	boll := indicator.Ind(lines).Boll()
	// 跨中线 双开
	if boll.CrossMB() {
		// 平掉开了的仓位
		if s.phasePositionExists() {
			log.Println("reset to close_all positions at", boll.CurrentPrice())
			//s.reset()
			s.closePhaseOrders()
		}

		//if s.opened {
		//	return nil
		//}

		if s.initLongOrder == nil {
			//longOrder, err := order.DualBuyLong(symbol, calcQty2(principal.SingleBetBalance(), boll.LastKline().Close))
			longOrder, err := order.DualBuyLong(s.symbol, principal.Qty())
			if err != nil {
				return err
			}
			log.Println("中线 long order", helper.ToJson(longOrder))
			s.initLongOrder = longOrder
			if err := s.storeLongOrder(longOrder); err != nil {
				log.Println(err)
			}
		}

		if s.initShortOrder == nil {
			//shortOrder, err := order.DualSellShort(symbol, calcQty2(principal.SingleBetBalance(), boll.LastKline().Close))
			shortOrder, err := order.DualSellShort(s.symbol, principal.Qty())
			if err != nil {
				return err
			}
			log.Println("中线 short order", helper.ToJson(shortOrder))
			s.initShortOrder = shortOrder
			if err := s.storeShortOrder(shortOrder); err != nil {
				log.Println(err)
			}
		}

		//s.opened = true
	} else if boll.IsFirstHalf() { // 上半段
		if err := s.phaseHandler(boll); err != nil {
			return err
		}
	} else if boll.IsSecondHalf() { // 下半段
		if err := s.phaseHandler(boll); err != nil {
			return err
		}
	} else if boll.CrossUP() { // 触碰上线
		if s.initLongOrder != nil {
			log.Println("上线 close")
			err := position.ClosePositionByOrderResp(s.initLongOrder)
			if err != nil {
				log.Println(err)
				//return err
			}
			s.initLongOrder = nil
		}
	} else if boll.CrossDN() { // 触碰下线
		if s.initShortOrder != nil {
			log.Println("下线 close")
			err := position.ClosePositionByOrderResp(s.initShortOrder)
			if err != nil {
				log.Println(err)
				//return err
			}
			s.initShortOrder = nil
		}
	}

	return nil
}

type profitResult struct {
	Profit   float64
	Ts       int64
	Datetime string
}

func newProfitResult(profit float64) profitResult {
	now := time.Now()
	return profitResult{
		Profit:   profit,
		Ts:       now.Unix(),
		Datetime: now.Format("2006-01-02 15:04:05"),
	}
}

func (s *Smooth) closePhasePositions() {
}

func (s *Smooth) closePhaseOrders() {
	err := position.CloseOrderResps(s.phaseOrders)
	if err != nil {
		log.Println(err)
		return
	}
	s.delPhaseKeys()
	s.phaseOrders = nil
}

func (s *Smooth) reset() {
	sum, err := position.CloseAllPositionsBySymbol(s.symbol)
	if err != nil {
		log.Println(err)
	}

	log.Println("Profit-sum:", sum)
	totalProfit := principal.ProfitSumUpdate(sum)
	log.Println("total Profit sum", totalProfit)
	b, _ := json.Marshal(newProfitResult(totalProfit))
	if err := cache.Client.LPush("profit_smooth", string(b)).Err(); err != nil {
		log.Println(err)
	}
	cache.Client.Del(prefix+"init_long", prefix+"init_short")

	s.opened = false
	s.initShortOrder = nil
	s.initLongOrder = nil
	s.phaseOrders = nil

	s.delPhaseKeys()
}

// 清除缓存 phase key
func (s *Smooth) delPhaseKeys() {
	var keys []string
	for i := 1; i <= s.p.segments; i++ {
		keys = append(keys, s.KeyPhase(i))
	}
	if err := cache.Client.Del(keys...); err != nil {
		log.Println(err)
		return
	}
}

func (s *Smooth) phaseHandler(boll indicator.Boll) error {
	phase := s.p.phase(boll)
	b, err := s.phaseExists(phase)
	if err != nil {
		return err
	}
	// phase不存在
	if !b {
		var o *futures.CreateOrderResponse
		//var err error
		if phase > 0 {
			o, err = order.DualSellShort(s.symbol, principal.Qty())
			if err == nil {
				s.phaseOrders = append(s.phaseOrders, o)
			}
		} else {
			o, err = order.DualBuyLong(s.symbol, principal.Qty())
			if err == nil {
				s.phaseOrders = append(s.phaseOrders, o)
			}
		}
		if err != nil {
			return err
		}
		log.Printf("半段开仓 phase: %v  order: %v\n", phase, helper.ToJson(o))
		if err := cache.Client.Set(s.KeyPhase(phase), 1, 0).Err(); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func (s *Smooth) phaseExists(phase int) (bool, error) {
	result, err := cache.Client.Exists(s.KeyPhase(phase)).Result()
	if err != nil {
		return false, err
	}
	if result > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *Smooth) phasePositionExists() bool {
	return len(cache.Client.Keys(pattern).Val()) > 0
}

// 监控 stop loss
func (s *Smooth) monitorLoss() {

}

func (s *Smooth) KeyPhase(state int) string {
	return prefix + strconv.Itoa(state)
}

func (s *Smooth) storeLongOrder(order *futures.CreateOrderResponse) error {
	orderLongBytes, err := json.Marshal(order)
	if err != nil {
		return err
	}
	if err := cache.Client.Set(prefix+"init_long", string(orderLongBytes), 0).Err(); err != nil {
		return err
	}
	return nil
}

func (s *Smooth) storeShortOrder(order *futures.CreateOrderResponse) error {
	orderShortBytes, err := json.Marshal(order)
	if err != nil {
		return err
	}
	if err := cache.Client.Set(prefix+"init_short", string(orderShortBytes), 0).Err(); err != nil {
		return err
	}
	return nil
}

//func storeInitOrders(orderLong *futures.CreateOrderResponse, orderShort *futures.CreateOrderResponse) error {
//	return nil
//}

func resumeInitOrder() (orderLong *futures.CreateOrderResponse, orderShort *futures.CreateOrderResponse, err error) {
	long, err := cache.Client.Get(prefix + "init_long").Result()
	if err != nil {
		return
	}

	short, err := cache.Client.Get(prefix + "init_short").Result()
	if err != nil {
		return nil, nil, err
	}

	if err = json.Unmarshal([]byte(long), orderLong); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(short), orderShort); err != nil {
		return
	}
	log.Printf("resume init long: %v short: %v\n", orderLong, orderShort)

	return
}
