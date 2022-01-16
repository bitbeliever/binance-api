package fapi

import (
	"encoding/json"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"os"
	"strings"
	"time"
)

// RealTimeKline 实时获取最新k线数据和实时计算boll
func RealTimeKline(symbol, interval string) {
	lines, err := KlineHistory(symbol, interval, 21, 0)
	if err != nil {
		log.Println(err)
		return
	}

	//for _, line := range lines {
	//	log.Println(time.UnixMilli(line.OpenTime).Format("15:04:05"), time.UnixMilli(line.CloseTime).Format("15:04:05"), line.Open, line.High, line.Low, line.Close)
	//}

	bRes := CalculateBollByFapiKline(lines)
	log.Println("from history:", toJson(bRes))

	ch := KlineStream(symbol, interval)
	//var pinfo positionInfo

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
			log.Println(toJson(bRes), lines[len(lines)-1].Close)

			// 穿过布林线
			//if isCrossingLine(bRes, lines[len(lines)-1]) {
			//	// todo
			//	log.Println("crossing... current market", toJson(bRes), toJson(lines[len(lines)-1]))
			//	//CreateOrder(symbol, futures.SideTypeBuy)
			//	return
			//}

			// 区分上下穿==================
			switch calCrossType(bRes, lines[len(lines)-1]) {
			case ascendCross:
				log.Println("asc crossing... current market", toJson(bRes), toJson(lines[len(lines)-1]))
				flog.Println("asc crossing... current market", toJson(bRes), toJson(lines[len(lines)-1]))
				CreateOrder(symbol, futures.SideTypeBuy, "0.05") // qty todo
				//pinfo.Lock()
				//pinfo.position += 0.05
				//pinfo.Unlock()
				go monitor(symbol, lines[len(lines)-1].Close, 10, futures.SideTypeSell)
			case descendCross:
				log.Println("desc crossing... current market", toJson(bRes), toJson(lines[len(lines)-1]))
				flog.Println("asc crossing... current market", toJson(bRes), toJson(lines[len(lines)-1]))
				CreateOrder(symbol, futures.SideTypeSell, "0.05") // qty todo
				//pinfo.Lock()
				//pinfo.position -= 0.05
				//pinfo.Unlock()
				go monitor(symbol, lines[len(lines)-1].Close, 10, futures.SideTypeBuy)
			}

			// ==================
		}
	}
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
func HistoryBoll(symbol string) {
	lines, err := KlineHistory(symbol, "15m", 22, 0)
	if err != nil {
		log.Println(err)
		return
	}
	lines = lines[:len(lines)-1]
	log.Println(time.UnixMilli(lines[len(lines)-1].OpenTime).Format("15:04:05"))
	log.Println(len(lines))

	log.Println(toJson(CalculateBollByFapiKline(lines)))
}
