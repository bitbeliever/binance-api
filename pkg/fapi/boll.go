package fapi

import (
	"github.com/adshao/go-binance/v2/futures"
	"math"
)

type bollResult struct {
	MB, UP, DN float64
}

var klineDataSum []futures.WsKline

func subscribeWsKline(ch chan futures.WsKline, out chan bollResult) {
	for {
		select {
		case kline := <-ch:
			klineDataSum = append(klineDataSum, kline)
			out <- calculateBoll(klineDataSum)
		}
	}
}

func calculateBoll(lines []futures.WsKline) bollResult {
	// N 时间
	N := len(lines)
	if N != 20 {
		// todo
	}

	var closeSum float64
	for _, line := range lines {
		closeSum += Str2Float64(line.Close)
	}
	MA := closeSum / float64(N)

	closeSum = 0
	for _, line := range lines {
		closeSum += math.Pow(Str2Float64(line.Close)-MA, 2)
	}
	// 标准差
	MD := math.Sqrt(closeSum / float64(N-1))

	return bollResult{
		MB: MA,
		UP: MA + MD*2,
		DN: MA - MD*2,
	}
}

func CalculateBollByFapiKline(lines []*futures.Kline) bollResult {
	// N 时间
	N := len(lines)

	var closeSum float64
	for _, line := range lines {
		closeSum += Str2Float64(line.Close)
	}
	MA := closeSum / float64(N)

	closeSum = 0
	for _, line := range lines {
		closeSum += math.Pow(Str2Float64(line.Close)-MA, 2)
	}
	// 标准差
	MD := math.Sqrt(closeSum / float64(N-1))

	return bollResult{
		MB: MA,
		UP: MA + MD*2,
		DN: MA - MD*2,
		//Time: time.UnixMilli(lines[N-1].StartTime),
	}
}
