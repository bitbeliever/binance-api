package fapi

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"math"
)

type bollResult struct {
	UP, MB, DN float64
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
		closeSum += helper.Str2Float64(line.Close)
	}
	MA := closeSum / float64(N)

	closeSum = 0
	for _, line := range lines {
		closeSum += math.Pow(helper.Str2Float64(line.Close)-MA, 2)
	}
	// 标准差
	MD := math.Sqrt(closeSum / float64(N-1))

	return bollResult{
		MB: MA,
		UP: MA + MD*2,
		DN: MA - MD*2,
	}
}

func isCrossingLine(bRes bollResult, line *futures.Kline) bool {
	MB := bRes.MB
	_ = MB
	UP := bRes.UP
	DN := bRes.DN

	open := helper.Str2Float64(line.Open)
	close := helper.Str2Float64(line.Close)

	// 穿过布林带上线
	if (open < UP && close > UP) ||
		(open > UP && close < UP) {
		//log.Println("crossed upper", toJson(line))
		return true
	} else if (open < DN && close > DN) ||
		(open > DN && close < DN) { // 穿过下线
		//log.Println("crossed down", toJson(line))
		return true
	}

	return false
}

// 计算穿针类型 todo 中线
func calCrossType(bRes bollResult, line *futures.Kline) crossType {
	UP := bRes.UP
	DN := bRes.DN

	open := helper.Str2Float64(line.Open)
	close := helper.Str2Float64(line.Close)

	// 上升 穿过布林带上线
	if open < UP && close > UP {
		return ascendCross
	} else if open > UP && close < UP {
		//log.Println("crossed upper", toJson(line))
		return descendCross
	} else if open < DN && close > DN { // 上升 穿过下线
		return ascendCross
	} else if open > DN && close < DN { // 下降 穿过下线
		//log.Println("crossed down", toJson(line))
		return descendCross
	}

	return noCross
}
