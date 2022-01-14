package fapi

import "github.com/adshao/go-binance/v2/futures"

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

// 计算穿针类型
func calCrossType(bRes bollResult, line *futures.Kline) crossType {
	UP := bRes.UP
	DN := bRes.DN

	open := Str2Float64(line.Open)
	close := Str2Float64(line.Close)

	// 下到上 穿过布林带上线
	if open < UP && close > UP {
		return ascendCross
	} else if open > UP && close < UP {
		//log.Println("crossed upper", toJson(line))
		return descendCross
	} else if open < DN && close > DN {
		return ascendCross
	} else if open > DN && close < DN { // 穿过下线
		//log.Println("crossed down", toJson(line))
		return descendCross
	}

	return noCross
}
