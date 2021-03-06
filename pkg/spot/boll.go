package spot

import (
	"github.com/adshao/go-binance/v2"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"math"
)

type bollResult struct {
	MB, UP, DN float64
}

var klineDataSum []binance.WsKline

func subscribeWsKline(ch chan binance.WsKline, out chan bollResult) {
	for {
		select {
		case kline := <-ch:
			klineDataSum = append(klineDataSum, kline)
			out <- calculateBoll(klineDataSum)
		}
	}
}

func calculateBoll(lines []binance.WsKline) bollResult {
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
	MD := math.Sqrt(closeSum / float64(N)) // binance using N instead of N-1

	return bollResult{
		MB: MA,
		UP: MA + MD*2,
		DN: MA - MD*2,
	}
}
