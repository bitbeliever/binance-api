package principal

import (
	"github.com/adshao/go-binance/v2/futures"
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
func StopPNL() float64 {
	tb.mu.RLock()
	defer tb.mu.RUnlock()

	//return 0.01
	return tb.balance * tb.stopRate
}

func GetBalance() float64 {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	return tb.balance
}

func SingleBetBalance() float64 {
	tb.mu.RLock()
	tb.mu.RUnlock()

	return tb.balance * tb.openPositionRate
}

func UpdateBalance(balance float64) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.balance = balance
}

func init() {
	if err := os.Remove("line.txt"); err != nil {
		log.Println(err)
	}

	balances, err := account.QueryBalance()
	must(err)
	tb = &totalBalance{}

	if len(balances) != 0 {
		// todo
		for _, balance := range balances {
			if balance.Asset == "USDT" {
				tb.balance = helper.Str2Float64(balance.CrossWalletBalance)
			}
		}
	}

	// todo
	tb.stopRate = 0.1
	tb.openPositionRate = 0.1
	log.Println("初始本金:", tb.balance)
	log.Println("stopPNL:", StopPNL())
	log.Println("principal", helper.ToJson(tb))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
