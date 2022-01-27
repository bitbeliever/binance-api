package spot

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

// Transfer 合约资金划转 POST /sapi/v1/futures/transfer  1: 现货账户向USDT合约账户划转 2: USDT合约账户向现货账户划转
func Transfer(asset, amount string) {
	resp, err := NewClient().NewFuturesTransferService().
		Type(binance.FuturesTransferTypeToFutures).
		Asset(asset).
		Amount(amount).
		Do(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(helper.ToJson(resp))
}
