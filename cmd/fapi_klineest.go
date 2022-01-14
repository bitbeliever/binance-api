package main

import (
	"github.com/bitbeliever/binance-api/pkg/fapi"
	"github.com/bitbeliever/binance-api/pkg/helper"
	"log"
	"time"
)

const (
	layout = "2006/01/02 15:04:05"
)

func main() {
	//fapi.KlineStreamTest()
	st, err := time.Parse(layout, "2022/01/11 19:00:00") // utc time
	if err != nil {
		panic(err)
	}
	log.Println(st)
	log.Println(time.UnixMilli(st.UnixMilli()))
	lines, err := fapi.KlineHistory(fapi.ETH, "15m", 21, st.UnixMilli())
	if err != nil {
		log.Println(err)
		return
	}

	//lines = lines[:len(lines)-1]
	for _, re := range lines {
		log.Println(time.UnixMilli(re.OpenTime).Format("15:04:05"), time.UnixMilli(re.CloseTime).Format("15:04:05"), re.Open, re.Close)
	}

	log.Println(len(lines))

	bRes := fapi.CalculateBollByFapiKline(lines)
	log.Println(helper.ToJson(bRes))

}
