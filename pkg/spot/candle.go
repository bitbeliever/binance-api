package spot

import (
	"context"
	"github.com/adshao/go-binance/v2"
)

func KlineHistory(symbol, interval string, limit int) ([]*binance.Kline, error) {
	klines, err := NewClient().NewKlinesService().Symbol(symbol).
		Interval(interval).Limit(limit).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return klines, nil
}
