package strategy

import (
	"fmt"
	"github.com/bitbeliever/binance-api/pkg/helper"
)

type Strategy interface {
	Do()
}

func calcQty(spend float64, closeStr string) string {
	price := helper.Str2Float64(closeStr)
	return fmt.Sprintf("%.2f", spend/price)
	//spend / price
	//return strconv.FormatFloat(math.Round(spend/price*100)/100, 'f', 10, 64)
}
