package fapi

import (
	"github.com/adshao/go-binance/v2/futures"
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
	//MD := math.Sqrt(closeSum / float64(N-1))
	MD := math.Sqrt(closeSum / float64(N)) // todo: 标准差: 除以N || N-1(币安公式)

	return bollResult{
		MB: MA,
		UP: MA + MD*2,
		DN: MA - MD*2,
		//Time: time.UnixMilli(lines[N-1].StartTime),
	}
}

func isCrossingLine(bRes bollResult, line *futures.Kline) bool {
	MB := bRes.MB
	_ = MB
	UP := bRes.UP
	DN := bRes.DN

	open := Str2Float64(line.Open)
	close := Str2Float64(line.Close)

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

	open := Str2Float64(line.Open)
	close := Str2Float64(line.Close)

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

// 穿过boll带中线
// !!todo +-1 for testing
func bollCrossMB(bRes bollResult, line *futures.Kline) bool {
	open := Str2Float64(line.Open)
	close := Str2Float64(line.Close)
	//return (open < bRes.MB && close >= bRes.MB) ||
	//	(open > bRes.MB && close <= bRes.MB)
	return (open < bRes.MB && (close+1) >= bRes.MB) ||
		(open > bRes.MB && (close-1) <= bRes.MB)
}

func bollCross(bRes bollResult, line *futures.Kline) bool {
	price := Str2Float64(line.Close)
	// !!todo for testing
	return (price >= bRes.UP-0.5) ||
		(price <= bRes.DN+0.5) ||
		//return (price >= bRes.UP) ||
		//	(price <= bRes.DN) ||
		bollCrossMB(bRes, line)
}
