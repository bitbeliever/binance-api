package strategy

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/cache"
	"github.com/bitbeliever/binance-api/pkg/fapi/indicator"
	"github.com/bitbeliever/binance-api/pkg/fapi/internal/principal"
	"github.com/bitbeliever/binance-api/pkg/fapi/order"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
	"strconv"
)

type SmoothTP struct {
	symbol string
	p      pyramid
}

func NewSmoothTP(symbol string) *SmoothTP {

	return &SmoothTP{
		symbol: symbol,
		//p:      newPyramid(false),
		p: pyramid{
			segments: 3,
		},
	}
}

func (s *SmoothTP) Do(lines []*futures.Kline) error {
	boll := indicator.Ind(lines).Boll()
	if boll.IsFirstHalf() || boll.IsSecondHalf() {
		if err := s.phaseHandler(boll); err != nil {
			return err
		}
	}

	return nil
}

func (s *SmoothTP) reset() {

}

func (s *SmoothTP) positionExists() bool {
	return len(cache.Client.Keys("smoothTP_*").Val()) > 0
}

func (s *SmoothTP) phaseHandler(boll indicator.Boll) error {
	phase := s.p.phase(boll)
	b, err := s.phaseExists(phase)
	if err != nil {
		return err
	}
	// phase不存在
	if !b {
		var o *futures.CreateOrderResponse
		var err error
		if phase > 0 {
			o, err = order.DualSellShortSL(s.symbol, principal.Qty(), boll.Result().MBStr())
		} else {
			o, err = order.DualBuyLongSL(s.symbol, principal.Qty(), boll.Result().MBStr())
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
func (s *SmoothTP) phaseExists(phase int) (bool, error) {
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

func (s *SmoothTP) KeyPhase(state int) string {
	return "smoothTP_" + strconv.Itoa(state)
}
