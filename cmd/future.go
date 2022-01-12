package main

import (
	"github.com/bitbeliever/binance-api/pkg/fur"
)

func main() {
	go fur.KlineStream()
	//go spotws.KlineStream()

	fur.AccountServiceTest()
	select {}
}
