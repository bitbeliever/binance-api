package configs

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"log"
)

var Cfg cfg

type cfg struct {
	Emu           bool
	Debug         bool
	KlineInterval string
	Key           key `toml:"key"`
	MaxLeverage   bool
	StopLimit     float64 // 止损
	TakeProfit    float64 // 止盈

	Strategy strategy
}

type key struct {
	ApiKey    string
	SecretKey string
}

type strategy struct {
	PyramidSegments int
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	_, err := toml.DecodeFile("configs/config.toml", &Cfg)
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(Cfg, "", "  ")
	if err != nil {
		panic(err)
	}
	log.Println("cfg string:", string(b))
}
