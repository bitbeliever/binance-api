package fapi

import (
	"encoding/json"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/fapi/indicator"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
	"os"
	"time"
)

// RealTimeKline 实时获取最新k线数据和实时计算boll
func RealTimeKline(symbol, interval string, limit int) (chan []*futures.Kline, error) {
	log.Printf("实时数据 symbol: %v, interval %v \n: ", symbol, interval)
	lines, err := KlineHistory(symbol, interval, limit)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	bRes := indicator.NewBoll(lines).Result()
	//writeToLineTestFile(lines)
	log.Println("from history:", helper.ToJson(bRes))

	ch := KlineStream(symbol, interval)
	//var s = strategy.NewDoubleOpenStrategy(symbol, lev)
	//var s = strategy.NewSmooth(symbol)
	linesCh := make(chan []*futures.Kline, 2<<10)

	go func() {
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
					//log.Println("next kline", time.UnixMilli(lastKline.CloseTime).Format("15:04:05"), time.UnixMilli(wsKline.StartTime).Format("15:04:05"), strings.Repeat("===", 50))
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
					lastKline = lines[len(lines)-1]

					//writeToLineTestFile(lines)
				}

				//bRes := indicator.NewBollResult(lines)
				//log.Println(toJson(bRes), lastKline.Close)

				linesCh <- lines
				//if err := s.Do(indicator.NewBoll(lines)); err != nil {
				//	log.Println(err)
				//	return
				//}
			}
		}
	}()
	return linesCh, nil
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

func writeToLineTestFile(lines []*futures.Kline) {
	f, err := os.OpenFile("line.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	_, err = f.Write(lineTimeData(lines))
	if err != nil {
		log.Println(err)
	}
}

// HistoryBoll 测试历史boll指标数据
func HistoryBoll(symbol string) {
	lines, err := KlineHistory(symbol, "15m", 22)
	if err != nil {
		log.Println(err)
		return
	}
	lines = lines[:len(lines)-1]
	log.Println(time.UnixMilli(lines[len(lines)-1].OpenTime).Format("15:04:05"))
	log.Println(len(lines))

	log.Println(helper.ToJson(indicator.NewBollResult(lines)))
}
