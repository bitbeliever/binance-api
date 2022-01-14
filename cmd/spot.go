package main

import (
	"github.com/bitbeliever/binance-api/pkg/spot"
	"log"
	"time"
)

func main() {
	spot.AccountService()

	spot.ServerTime()
	log.Println("now", time.Now())
}
