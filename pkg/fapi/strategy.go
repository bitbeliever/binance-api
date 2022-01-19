package fapi

import (
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"log"
)

// V1 穿过布林线
func v1Cross(symbol string, bRes bollResult, lastKline *futures.Kline) error {
	if isCrossingLine(bRes, lastKline) {
		// todo
		log.Println("crossing... current market", toJson(bRes), toJson(lastKline))
		//CreateOrder(symbol, futures.SideTypeBuy)
		return nil
	}

	return nil
}

// V2 区分上下穿==================
func distUpperLowerCross(symbol string, bRes bollResult, lastKline *futures.Kline) error {

	switch calCrossType(bRes, lastKline) {
	case ascendCross:
		log.Println("asc crossing... current market", toJson(bRes), toJson(lastKline))
		flog.Println("asc crossing... current market", toJson(bRes), toJson(lastKline))
		CreateOrder(symbol, futures.SideTypeBuy, "0.05") // qty todo
		//pinfo.Lock()
		//pinfo.position += 0.05
		//pinfo.Unlock()
		go monitor(symbol, lastKline.Close, 10, futures.SideTypeSell)
	case descendCross:
		log.Println("desc crossing... current market", toJson(bRes), toJson(lastKline))
		flog.Println("asc crossing... current market", toJson(bRes), toJson(lastKline))
		CreateOrder(symbol, futures.SideTypeSell, "0.05") // qty todo
		//pinfo.Lock()
		//pinfo.position -= 0.05
		//pinfo.Unlock()
		go monitor(symbol, lastKline.Close, 10, futures.SideTypeBuy)
	}
	return nil
}

// v3 中线双开策略
// V3 穿过中线, 双开 有未平的仓========================================
// todo yesterday: 达到中线后, 空/多判断, 开其中不存在的(可能多单或者空单只存其一)
// todo 本金10%
func mbDoubleOpenPosition(symbol string, bRes bollResult, lastKline *futures.Kline) error {
	if !bollCross(bRes, lastKline) {
		return nil
	}

	//log.Println("crossing:", toJson(bRes), lastKline.Close)

	// 止损已经挂起, 直接返回
	if principal.stopLossFn != nil {
		return nil
	}

	// !!todo query too many times cause error
	positions, err := QueryAccountPositions()
	if err != nil {
		log.Println(err)
		return err
	}

	if bollCrossMB(bRes, lastKline) {
		// 有仓位, 不开新仓 todo
		if len(positions) != 0 {
			//log.Println("positions exists")
			return nil
		}
		// open position
		log.Println("达到中线======== to open position", toJson(bRes), lastKline.Close)
		// buy long
		orderBuy, err := CreateOrderDual(symbol, futures.SideTypeBuy, futures.PositionSideTypeLong, calcQty(10, lastKline.Close))
		if err != nil {
			log.Println(err)
			return err
		}
		// short sell
		orderSell, err := CreateOrderDual(symbol, futures.SideTypeSell, futures.PositionSideTypeShort, calcQty(10, lastKline.Close))
		if err != nil { // todo close
			log.Println(err)
			return err
		}

		log.Println("buy order", toJson(orderBuy))
		log.Println("sell order", toJson(orderSell))
		principal.orderLong = orderBuy
		principal.orderShort = orderSell
	} else if Str2Float64(lastKline.Close) >= bRes.UP { // 触碰上线 平多单 止盈
		if len(positions) == 0 {
			return nil
		}

		// 止盈
		// 还未平仓
		if principal.orderLong != nil {
			if err := closePosition(principal.orderLong); err != nil {
				return err
			}
			principal.orderLong = nil
		}
		// 止损 todo
		if err := closePosition(principal.orderShort); err != nil {
			return err
		}

	} else if Str2Float64(lastKline.Close) <= bRes.DN { // 触碰下线
		if len(positions) == 0 {
			return nil
		}

		// 止盈
		// 还未平仓
		if principal.orderShort != nil {
			if err := closePosition(principal.orderShort); err != nil {
				return err
			}
			principal.orderShort = nil
		}
		// 止损todo
		if err := closePosition(principal.orderLong); err != nil {
			return err
		}

	}

	return nil
}

func calcQty(spend float64, closeStr string) string {
	price := Str2Float64(closeStr)
	return fmt.Sprintf("%.2f", spend/price)
	//spend / price
	//return strconv.FormatFloat(math.Round(spend/price*100)/100, 'f', 10, 64)
}

func Test() {
	log.Println(calcQty(10, "3136.32"))
	log.Println(10 / 3136.32)
}
