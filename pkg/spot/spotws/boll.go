package spotws

import (
	"github.com/adshao/go-binance/v2"
	"math"
	"time"
)

type bollResult struct {
	MB, UP, DN float64
	Time       time.Time
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

// todo 首次读取, N = 1
func calculateBoll(lines []binance.WsKline) bollResult {
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
	var MD float64 // todo
	if N == 1 {
		MD = 0
	} else {
		MD = math.Sqrt(closeSum / float64(N-1))
	}
	//MD := math.Sqrt(closeSum / float64(N-1))

	return bollResult{
		MB:   MA,
		UP:   MA + MD*2,
		DN:   MA - MD*2,
		Time: time.UnixMilli(lines[N-1].StartTime),
	}
}
