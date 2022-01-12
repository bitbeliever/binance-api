package main

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/spot"
)

type A struct {
	A int
	B int
}

func main() {
	ch := make(chan *futures.WsUserDataEvent, 2^10)
	go spot.UserDataStream()
	go fapi.UserDataStream(ch)

	select {}
}
