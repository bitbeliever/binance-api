package main

import "github.com/bitbeliever/binance-api/pkg/fapi"

const (
	chBuf = 2 << 10 // (2 ^ 11)
)

func main() {
	fapi.RecvUserDataStream()
}
