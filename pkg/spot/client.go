package spot

import (
	"github.com/adshao/go-binance/v2"
	"github.com/bitbeliever/binance-api/configs"
)

type Client struct {
	*binance.Client
}

func NewClient() Client {
	return Client{
		binance.NewClient(configs.Cfg.Key.ApiKey, configs.Cfg.Key.SecretKey),
	}
}
