package indicator

import "github.com/adshao/go-binance/v2/futures"

type Indicator interface {
	Price() float64
	CurrentPrice() float64
}

type Ind []*futures.Kline

func (ind Ind) Ma() Ma {
	return NewMa(ind)
}

func (ind Ind) Boll() Boll {
	return NewBoll(ind)
}
