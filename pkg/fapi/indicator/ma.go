package indicator

import (
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/helper"
)

type Ma struct {
	lines []*futures.Kline
}

func NewMa(lines []*futures.Kline) Ma {
	return Ma{
		lines: lines,
	}
}

func (m Ma) AveragePrice() string {
	var sum float64
	for _, v := range m.lines {
		sum += helper.Str2Float64(v.Close)
	}
	return fmt.Sprintf("%.2f", sum/float64(len(m.lines)))
}

func (m Ma) CurrentPrice() string {
	return m.lines[len(m.lines)-1].Close
}
