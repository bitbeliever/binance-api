package fapi

import (
	"encoding/json"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"os"
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
			// 同一周期, 更新最后一个candle
			if wsKline.StartTime == lastKline.OpenTime && wsKline.EndTime == lastKline.CloseTime {
				lastKline.Low = wsKline.Low
				lastKline.High = wsKline.High
				lastKline.Open = wsKline.Open
				lastKline.Close = wsKline.Close
				//lines[len(lines)-1] = lastKline
			} else { // new cycle
				log.Println("next kline", time.UnixMilli(lastKline.CloseTime).Format("15:04:05"), time.UnixMilli(wsKline.StartTime).Format("15:04:05"), strings.Repeat("===", 50))
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

				/////////// test
				f, err := os.OpenFile("line.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				_, err = f.Write(lineTimeData(lines))
				if err != nil {
					log.Println(err)
					return
				}
				/////////// test
			}

			//log.Println("sub", wsKline.StartTime-lastKline.CloseTime)
			//log.Println("wrong kline", time.UnixMilli(lastKline.OpenTime).Format("15:04:05"), time.UnixMilli(lastKline.CloseTime).Format("15:04:05"),
			//	time.UnixMilli(wsKline.StartTime).Format("15:04:05"), time.UnixMilli(wsKline.EndTime).Format("15:04:05"))
			//continue
			bRes = CalculateBollByFapiKline(lines)
			log.Println(toJson(bRes))

			// 穿过布林线
			if isCrossingLine(bRes, lines[len(lines)-1]) {
				// todo
				//CreateOrder()
			}
		}
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

// lineTimeData 测试实时k线数据辅助
func lineTimeData(lines []*futures.Kline) []byte {
	type d struct {
		StartTime string
		EndTime   string
	}
	var data []d
	for _, line := range lines {
		data = append(data, d{
			StartTime: time.UnixMilli(line.OpenTime).Format("15:04:05"),
			EndTime:   time.UnixMilli(line.CloseTime).Format("15:04:05"),
		})
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
		return nil
	}
	b = append(b, '\n')
	return b
}

// HistoryBoll 测试历史boll指标数据
func HistoryBoll() {
	lines, err := KlineHistory("ETHUSDT", "15m", 22, 0)
	if err != nil {
		log.Println(err)
		return
	}
	lines = lines[:len(lines)-1]
	log.Println(time.UnixMilli(lines[len(lines)-1].OpenTime).Format("15:04:05"))
	log.Println(len(lines))

	log.Println(toJson(CalculateBollByFapiKline(lines)))
}
