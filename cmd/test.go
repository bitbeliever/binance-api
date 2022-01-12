package main

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/fur"
	"github.com/bitbeliever/binance-api/pkg/spot"
	"reflect"
)

type A struct {
	A int
	B int
}

func main() {
	var a = A{1, 2}
	var b = struct {
		A, B int
	}{1, 2}
	_ = b
	var c = A{1, 2}
	_ = c

	println(reflect.DeepEqual(a, b))

	ch := make(chan *futures.WsUserDataEvent, 2^10)
	go spot.UserDataStream()
	go fur.UserDataStream(ch)

	select {}
}
