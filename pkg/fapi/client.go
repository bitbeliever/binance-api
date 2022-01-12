package fapi

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/configs"
)

type Client struct {
	*futures.Client
}

func NewClient() Client {
	return Client{
		futures.NewClient(configs.Cfg.Key.ApiKey, configs.Cfg.Key.SecretKey),
	}
}
