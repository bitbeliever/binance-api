package fapi

import (
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"strings"
	"time"
)

// RealTimeKline 实时计算boll
func RealTimeKline() {
	lines, err := KlineHistory("ETHUSDT", "15m", 21, 0)
	if err != nil {
		log.Println(err)
		return
	}

	for _, line := range lines {
		log.Println(time.UnixMilli(line.OpenTime).Format("15:04:05"), time.UnixMilli(line.CloseTime).Format("15:04:05"), line.Open, line.High, line.Low, line.Close)
	}

	bRes := CalculateBollByFapiKline(lines)
	log.Println("from history:", toJson(bRes))

	ch := KlineStream("ETHUSDT", "15m")

	for {
		select {
		case wsKline := <-ch:
			lastKline := lines[len(lines)-1]
			if wsKline.StartTime == lastKline.OpenTime && wsKline.EndTime == lastKline.CloseTime {
				lastKline.Low = wsKline.Low
				lastKline.High = wsKline.High
				lastKline.Open = wsKline.Open
				lastKline.Close = wsKline.Close
				//lines[len(lines)-1] = lastKline
			} else { // new cycle
				log.Println("next kline", time.UnixMilli(lastKline.CloseTime).Format("15:04:05"), time.UnixMilli(wsKline.StartTime).Format("15:04:05"), strings.Repeat("===", 10))
				lines = lines[1:]
				lines = append(lines, &futures.Kline{
					OpenTime:  wsKline.StartTime,
					Open:      wsKline.Open,
					High:      wsKline.High,
					Low:       wsKline.Low,
					Close:     wsKline.Close,
					Volume:    wsKline.Volume,
					CloseTime: wsKline.EndTime,
				})
			}

			//log.Println("sub", wsKline.StartTime-lastKline.CloseTime)
			//log.Println("wrong kline", time.UnixMilli(lastKline.OpenTime).Format("15:04:05"), time.UnixMilli(lastKline.CloseTime).Format("15:04:05"),
			//	time.UnixMilli(wsKline.StartTime).Format("15:04:05"), time.UnixMilli(wsKline.EndTime).Format("15:04:05"))
			//continue
			bRes = CalculateBollByFapiKline(lines)
			log.Println(toJson(bRes))

			if isCrossingLine(bRes, lines[len(lines)-1]) {
				// todo
				//CreateOrder()
			}
		}
	}
}

func isCrossingLine(bRes bollResult, line *futures.Kline) bool {
	MB := bRes.MB
	UP := bRes.UP
	DN := bRes.DN

	open := Str2Float64(line.Open)
	close := Str2Float64(line.Close)

	// 上涨, 穿过布林带上线
	if line.Close > line.Open && open < UP && close > UP {
		log.Println("crossed upper", toJson(line))
		return true
	} else if open < MB && close > MB { // 上涨 穿过中线
		log.Println("crossed middle", toJson(line))
		return true
	} else if open < DN && close > DN {
		log.Println("crossed down", toJson(line))
		return true
	}

	return false
}
