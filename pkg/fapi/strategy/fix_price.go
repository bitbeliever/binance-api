package strategy

import (
	"github.com/bitbeliever/binance-api/pkg/fapi/indicator"
)

type FixPrice struct {
	symbol string
	price  float64
}

func NewFixPrice(symbol string, price float64) *FixPrice {
	return &FixPrice{
		symbol: symbol,
		price:  price,
	}
}

func (f FixPrice) Do(ma indicator.Ma) error {
	//if ma.Price() >= f.price {
	//	position.CloseAllPositions()
	//}
	return nil
}
