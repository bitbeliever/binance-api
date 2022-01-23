package fapi

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/fapi/client"
)

// KlineHistory 获取历史k线
func KlineHistory(symbol string, interval string, limit int) ([]*futures.Kline, error) {
	res, err := client.NewClient().NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(limit).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	return res, nil
}
