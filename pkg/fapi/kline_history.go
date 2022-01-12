package fapi

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
)

// KlineHistory 获取历史k线
func KlineHistory(symbol string, interval string, limit int, startTime int64) ([]*futures.Kline, error) {
	res, err := NewClient().NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		//EndTime(startTime).
		Limit(limit).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	return res, nil
}
