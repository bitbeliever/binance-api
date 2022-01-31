package strategy

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/fapi/indicator"
	"log"
)

type Show struct{}

func NewShow() Show {
	return Show{}
}

func (s Show) Do(lines []*futures.Kline) error {
	ma := indicator.Ind(lines).Ma()
	log.Println(ma.AveragePrice())
	return nil
}
