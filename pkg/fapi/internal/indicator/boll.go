package indicator

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"math"
)

type Boll struct {
	klines []*futures.Kline
	res    BollResult
}
type BollResult struct {
	UP, MB, DN float64
}

func NewBoll(klines []*futures.Kline) Boll {
	b := Boll{
		klines: klines,
	}
	b.res = b.calculate()

	return b
}

func (b Boll) calculate() BollResult {
	lines := b.klines
	// N 时间
	N := len(lines)

	var closeSum float64
	for _, line := range lines {
		closeSum += helper.Str2Float64(line.Close)
	}
	MA := closeSum / float64(N)

	closeSum = 0
	for _, line := range lines {
		closeSum += math.Pow(helper.Str2Float64(line.Close)-MA, 2)
	}
	// 标准差
	//MD := math.Sqrt(closeSum / float64(N-1))
	MD := math.Sqrt(closeSum / float64(N)) // 标准差: 除以N || N-1(币安公式)

	b.res = BollResult{
		MB: MA,
		UP: MA + MD*2,
		DN: MA - MD*2,
	}

	return b.res
}

// CrossMB 穿过boll带中线
func (b Boll) CrossMB() bool {
	line := b.klines[len(b.klines)-1]
	open := helper.Str2Float64(line.Open)
	close := helper.Str2Float64(line.Close)
	return (open < b.res.MB && close >= b.res.MB) ||
		(open > b.res.MB && close <= b.res.MB)
	//return (open < b.res.MB && (close+1) >= b.res.MB) ||
	//	(open > b.res.MB && (close-1) <= res.MB)
}

// Cross 布林穿过线
func (b Boll) Cross() bool {
	bRes := b.res
	line := b.klines[len(b.klines)-1]

	price := helper.Str2Float64(line.Close)

	//return (price >= res.UP-0.5) ||
	//	(price <= res.DN+0.5) ||
	return (price >= bRes.UP) ||
		(price <= bRes.DN) ||
		b.CrossMB()
}

func (b Boll) CrossUP() bool {
	return helper.Str2Float64(b.LastKline().Close) >= b.res.UP
}

func (b Boll) CrossDN() bool {
	return helper.Str2Float64(b.LastKline().Close) <= b.res.DN
}

func (b Boll) LastKline() *futures.Kline {
	return b.klines[len(b.klines)-1]
}

func (b Boll) Result() BollResult {
	return b.res
}

func NewBollResult(klines []*futures.Kline) BollResult {
	return NewBoll(klines).calculate()
}
