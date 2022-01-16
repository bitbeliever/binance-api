package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
)

func main() {

	mode, err := fapi.PositionMode()
	log.Println(helper.ToJson(mode), err)
}
