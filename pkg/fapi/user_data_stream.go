package fapi

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/configs"
	"log"
	"time"
)

/*
UserDataStream
账户更新 account_update
{
  "e": "ACCOUNT_UPDATE",                // 事件类型
  "E": 1564745798939,                   // 事件时间
  "T": 1564745798938 ,                  // 撮合时间
  "a":                                  // 账户更新事件
    {
      "m":"ORDER",                      // 事件推出原因
      "B":[                             // 余额信息
        {
          "a":"USDT",                   // 资产名称
          "wb":"122624.12345678",       // 钱包余额
          "cw":"100.12345678",          // 除去逐仓仓位保证金的钱包余额
          "bc":"50.12345678"            // 除去盈亏与交易手续费以外的钱包余额改变量
        },
      ],
      "P":[
       {
          "s":"BTCUSDT",            // 交易对
          "pa":"0",                 // 仓位
          "ep":"0.00000",            // 入仓价格
          "cr":"200",               // (费前)累计实现损益
          "up":"0",                     // 持仓未实现盈亏
          "mt":"isolated",              // 保证金模式
          "iw":"0.00000000",            // 若为逐仓，仓位保证金
          "ps":"BOTH"                   // 持仓方向
       }，
      ]
    }
}

新的交易 order_trade_update
{
  "e":"ORDER_TRADE_UPDATE",         // 事件类型
  "E":1568879465651,                // 事件时间
  "T":1568879465650,                // 撮合时间
  "o":{
    "s":"BTCUSDT",                  // 交易对
    "c":"TEST",                     // 客户端自定订单ID
      // 特殊的自定义订单ID:
      // "autoclose-"开头的字符串: 系统强平订单
      // "adl_autoclose": ADL自动减仓订单
    "S":"SELL",                     // 订单方向
    "o":"TRAILING_STOP_MARKET", // 订单类型
    "f":"GTC",                      // 有效方式
    "q":"0.001",                    // 订单原始数量
    "p":"0",                        // 订单原始价格
    "ap":"0",                       // 订单平均价格
    "sp":"7103.04",                 // 条件订单触发价格，对追踪止损单无效
    "x":"NEW",                      // 本次事件的具体执行类型
    "X":"NEW",                      // 订单的当前状态
    "i":8886774,                    // 订单ID
    "l":"0",                        // 订单末次成交量
    "z":"0",                        // 订单累计已成交量
    "L":"0",                        // 订单末次成交价格
    "N": "USDT",                    // 手续费资产类型
    "n": "0",                       // 手续费数量
    "T":1568879465651,              // 成交时间
    "t":0,                          // 成交ID
    "b":"0",                        // 买单净值
    "a":"9.91",                     // 卖单净值
    "m": false,                     // 该成交是作为挂单成交吗？
    "R":false   ,                   // 是否是只减仓单
    "wt": "CONTRACT_PRICE",         // 触发价类型
    "ot": "TRAILING_STOP_MARKET",   // 原始订单类型
    "ps":"LONG"                     // 持仓方向
    "cp":false,                     // 是否为触发平仓单; 仅在条件订单情况下会推送此字段
    "AP":"7476.89",                 // 追踪止损激活价格, 仅在追踪止损单时会推送此字段
    "cr":"5.0",                     // 追踪止损回调比例, 仅在追踪止损单时会推送此字段
    "rp":"0"                       // 该交易实现盈亏
  }

}
*/
func UserDataStream(ch chan *futures.WsUserDataEvent) {
	c := futures.NewClient(configs.Cfg.Key.ApiKey, configs.Cfg.Key.SecretKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	listenKey, err := c.NewStartUserStreamService().Do(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("fapi key:", listenKey)

	doneCh, stopCh, err := futures.WsUserDataServe(listenKey, func(event *futures.WsUserDataEvent) {
		ch <- event
	}, func(err error) {
		log.Println("fapi data serve err", err)
		return
	})

	if err != nil {
		log.Println(err)
		return
	}
	_ = doneCh
	_ = stopCh

	go keepListenKeyAlive(c, listenKey)
}

// keep alive
func keepListenKeyAlive(client *futures.Client, listenKey string) {

	ticker := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-ticker.C:
			//log.Println("user data stream keeping alive")
			if err := client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(context.Background()); err != nil {
				log.Println(err)
			}
		}
	}
}

func UserDataStreamTest() {
	//ws.KlineStream()
	ch := make(chan *futures.WsUserDataEvent, chBuf)
	UserDataStream(ch)

	for {
		select {
		case event := <-ch:
			formatPrintEvent(event)
		}
	}
}

// PnL: Profit and Loss
func formatPrintEvent(event *futures.WsUserDataEvent) {
	// 账户更新事件 account update
	if event.Event == futures.UserDataEventTypeAccountUpdate {
		log.Printf("事件 %v, Time %v TranTime %v\n", event.Event, time.UnixMilli(event.Time).Format(layout), time.UnixMilli(event.TransactionTime).Format(layout))
		//log.Printf("Reason %v\n", event.AccountUpdate.Reason)
		for _, balance := range event.AccountUpdate.Balances {
			log.Printf("资产 %v, 余额 %v, 除去逐仓仓位保证金的钱包余额 %v, 该变量 %v\n", balance.Asset, balance.Balance, balance.CrossWalletBalance, balance.ChangeBalance)
		}
		for _, position := range event.AccountUpdate.Positions {
			log.Printf("交易对 %v, 仓位 %v, 方向 %v, 费前累计损益 %v, 持仓未实现盈亏 %v, 入仓价格 %v, 保证金模式 %v, 逐仓保证金: %v\n", position.Symbol, position.Amount, position.Side, position.AccumulatedRealized, position.UnrealizedPnL, position.EntryPrice, position.MarginType, position.IsolatedWallet)
		}
		log.Println("================================================================================")
	} else if event.Event == futures.UserDataEventTypeOrderTradeUpdate {
		// 订单更新事件 order trade update
		log.Printf("事件 %v, Time %v TranTime %v\n", event.Event, time.UnixMilli(event.Time).Format(layout), time.UnixMilli(event.TransactionTime).Format(layout))
		log.Printf("交易对 %v, pnl %v, tradeTime %v \n", event.OrderTradeUpdate.Symbol, event.OrderTradeUpdate.RealizedPnL, time.UnixMilli(event.OrderTradeUpdate.TradeTime))
		log.Printf("状态 %v, 执行类型 %v, 有效方式 %v \n", event.OrderTradeUpdate.Status, event.OrderTradeUpdate.ExecutionType, event.OrderTradeUpdate.TimeInForce)
		log.Printf("tradeID %v, 类型 %v, 方向 %v \n", event.OrderTradeUpdate.TradeID, event.OrderTradeUpdate.Type, event.OrderTradeUpdate.Side)
		log.Printf("原始数量: %v, 原始价格: %v, 平均价格: %v \n", event.OrderTradeUpdate.OriginalQty, event.OrderTradeUpdate.OriginalPrice, event.OrderTradeUpdate.AveragePrice)
		log.Printf("末次成交量: %v, 末次成交价格: %v, 累计成交量: %v \n", event.OrderTradeUpdate.LastFilledQty, event.OrderTradeUpdate.LastFilledPrice, event.OrderTradeUpdate.AccumulatedFilledQty)
		log.Printf("买单净值 %v, 卖单净值 %v, 手续费数量 %v \n", event.OrderTradeUpdate.BidsNotional, event.OrderTradeUpdate.AsksNotional, event.OrderTradeUpdate.Commission)
		log.Println("================================================================================")
	} else {
		// 其他事件
		log.Println("other event occurs", toJson(event))
	}
}
