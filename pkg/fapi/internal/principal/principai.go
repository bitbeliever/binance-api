package principal

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/configs"
	"github.com/bitbeliever/binance-api/pkg/account"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
	"os"
	"sync"
)

var (
	//principal *totalBalance
	tb *totalBalance
)

type totalBalance struct {
	mu sync.RWMutex
	// usdt
	Balance          float64
	OpenPositionRate float64 // 单次下单率 0.1

	// 双开
	orderLong, orderShort,
	closeLongOrder, closeShortOrder *futures.CreateOrderResponse

	StopRate   float64 // 止损 0.1 for now
	TakeProfit float64 // 止盈率

	stopLossFn func()
	profitSum  float64
}

// 止损的pnl
func StopPNL() float64 {
	tb.mu.RLock()
	defer tb.mu.RUnlock()

	//return 0.01
	return tb.Balance * tb.StopRate
}

func GetBalance() float64 {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	return tb.Balance
}

func SingleBetBalance() float64 {
	tb.mu.RLock()
	tb.mu.RUnlock()

	return tb.Balance * tb.OpenPositionRate
}

func UpdateBalance(balance float64) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.Balance = balance
}

func ProfitSumUpdate(profit float64) float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.profitSum += profit
	return tb.profitSum
}

// Qty !!todo
func Qty() string {
	return configs.Cfg.Qty
}

func init() {
	if err := os.Remove("line.txt"); err != nil {
		log.Println(err)
	}

	balances, err := account.QueryBalance()
	must(err)
	tb = &totalBalance{}

	if len(balances) != 0 {
		// 只计算USDT的余额
		for _, balance := range balances {
			if balance.Asset == "USDT" {
				tb.Balance = helper.Str2Float64(balance.CrossWalletBalance)
			}
		}
	}

	// todo
	tb.StopRate = 0.1
	tb.OpenPositionRate = 0.1
	log.Println("初始本金:", tb.Balance)
	log.Println("stopPNL:", StopPNL())
	log.Println("principal", helper.ToJson(tb))
	log.Println("principal singleBetBalance()", SingleBetBalance())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
