package configs

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"io"
	"log"
	"os"
	"time"
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
	Qty           string
	Symbol        string

	Strategy strategy
	Redis    redis
}

type key struct {
	ApiKey    string
	SecretKey string
}

type redis struct {
	Host     string
	Port     int
	Password string
}

type strategy struct {
	PyramidSegments int
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	f, err := os.Create(time.Now().Format("log/2006-01-02_15_04_05") + ".log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	_, err = toml.DecodeFile("configs/config.toml", &Cfg)
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(Cfg, "", "  ")
	if err != nil {
		panic(err)
	}
	log.Println("cfg string:", string(b))
}
