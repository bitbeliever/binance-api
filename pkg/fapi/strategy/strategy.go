package strategy

import (
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bitbeliever/binance-api/pkg/helper"
)

type Strategy interface {
	Do([]*futures.Kline) error
}

//func calcQty(spend float64, closeStr string) string {
//	price := helper.Str2Float64(closeStr)
//	return fmt.Sprintf("%.2f", spend/price)
//	//return strconv.FormatFloat(math.Round(spend/price*100)/100, 'f', 10, 64)
//}

// 花费的本金 * 杠杆 / 当前价格
func calcQty(spend float64, closeStr string, leverage int) string {
	price := helper.Str2Float64(closeStr)
	return fmt.Sprintf("%.2f", spend*float64(leverage)/price)
	//return strconv.FormatFloat(math.Round(spend/price*100)/100, 'f', 10, 64)
}

func calcQty2(spend float64, closeStr string) string {
	price := helper.Str2Float64(closeStr)
	return fmt.Sprintf("%.2f", spend/price)
}
