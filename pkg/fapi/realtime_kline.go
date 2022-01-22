package fapi

import (
	"encoding/json"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"os"
	"sync"
	"time"
)

var (
	principal *totalBalance
)

type totalBalance struct {
	mu sync.RWMutex
	// usdt
	balance          float64
	openPositionRate float64 // 单次下单率 0.1

	// 双开
	orderLong, orderShort,
	closeLongOrder, closeShortOrder *futures.CreateOrderResponse

	stopRate   float64 // 止损 0.1 for now
	takeProfit float64 // 止盈率

	stopLossFn func()
}

// 止损的pnl
func (tb *totalBalance) stopPNL() float64 {
	tb.mu.RLock()
	defer tb.mu.RUnlock()

	return 0.01
	return tb.balance * tb.stopRate
}

func (tb *totalBalance) getBalance() float64 {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	return tb.balance
}

func (tb *totalBalance) singleBetBalance() float64 {
	tb.mu.RLock()
	tb.mu.RUnlock()

	return tb.balance * tb.openPositionRate
}

func (tb *totalBalance) updateBalance(balance float64) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.balance = balance
}

func init() {
	if err := os.Remove("line.txt"); err != nil {
		log.Println(err)
	}

	balances, err := QueryBalance()
	must(err)
	principal = &totalBalance{}

	if len(balances) != 0 {
		// todo
		for _, balance := range balances {
			if balance.Asset == "USDT" {
				principal.balance = Str2Float64(balance.CrossWalletBalance)
			}
		}
	}

	// todo
	principal.stopRate = 0.1
	principal.openPositionRate = 0.1
	log.Println("初始本金:", principal.balance)
	log.Println("stopPNL:", principal.stopPNL())
	log.Println("principal", principal)
}

// RealTimeKline 实时获取最新k线数据和实时计算boll
func RealTimeKline(symbol, interval string) {
	log.Printf("实时数据 symbol: %v, interval %v \n: ", symbol, interval)
	lines, err := KlineHistory(symbol, interval, 21, 0)
	if err != nil {
		log.Println(err)
		return
	}

	//for _, line := range lines {
	//	log.Println(time.UnixMilli(line.OpenTime).Format("15:04:05"), time.UnixMilli(line.CloseTime).Format("15:04:05"), line.Open, line.High, line.Low, line.Close)
	//}

	bRes := CalculateBollByFapiKline(lines)
	writeToLineTestFile(lines)
	log.Println("from history:", toJson(bRes))

	ch := KlineStream(symbol, interval)
	var s = newDoubleOpenStrategy()
	// 设置止盈 call only once
	go s.monitorOrderTP(s.subscribeUpper(), s.subscribeLower())

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

				writeToLineTestFile(lines)
			}

			bRes = CalculateBollByFapiKline(lines)
			//log.Println(toJson(bRes), lastKline.Close)

			// v3 double open position
			//if err := s.mbDoubleOpenPosition(symbol, bRes, lastKline); err != nil {
			if err := s.mbDoubleOpenPositionByChannel(symbol, bRes, lastKline); err != nil {
				log.Println(err)
				return
			}
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

func writeToLineTestFile(lines []*futures.Kline) {
	f, err := os.OpenFile("line.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	_, err = f.Write(lineTimeData(lines))
	if err != nil {
		log.Println(err)
	}
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
