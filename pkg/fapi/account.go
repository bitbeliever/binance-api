package fapi

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"time"
)

// QueryAccount 账户查询 账户信息 fapi/v1/account
/*
{
    "feeTier": 0,  // 手续费等级
    "canTrade": true,  // 是否可以交易
    "canDeposit": true,  // 是否可以入金
    "canWithdraw": true, // 是否可以出金
    "updateTime": 0,
    "totalInitialMargin": "0.00000000",  // 但前所需起始保证金总额(存在逐仓请忽略), 仅计算usdt资产
    "totalMaintMargin": "0.00000000",  // 维持保证金总额, 仅计算usdt资产
    "totalWalletBalance": "23.72469206",   // 账户总余额, 仅计算usdt资产
    "totalUnrealizedProfit": "0.00000000",  // 持仓未实现盈亏总额, 仅计算usdt资产
    "totalMarginBalance": "23.72469206",  // 保证金总余额, 仅计算usdt资产
    "totalPositionInitialMargin": "0.00000000",  // 持仓所需起始保证金(基于最新标记价格), 仅计算usdt资产
    "totalOpenOrderInitialMargin": "0.00000000",  // 当前挂单所需起始保证金(基于最新标记价格), 仅计算usdt资产
    "totalCrossWalletBalance": "23.72469206",  // 全仓账户余额, 仅计算usdt资产
    "totalCrossUnPnl": "0.00000000",    // 全仓持仓未实现盈亏总额, 仅计算usdt资产
    "availableBalance": "23.72469206",       // 可用余额, 仅计算usdt资产
    "maxWithdrawAmount": "23.72469206"     // 最大可转出余额, 仅计算usdt资产
    "assets": [
        {
            "asset": "USDT",        //资产
            "walletBalance": "23.72469206",  //余额
            "unrealizedProfit": "0.00000000",  // 未实现盈亏
            "marginBalance": "23.72469206",  // 保证金余额
            "maintMargin": "0.00000000",    // 维持保证金
            "initialMargin": "0.00000000",  // 当前所需起始保证金
            "positionInitialMargin": "0.00000000",  // 持仓所需起始保证金(基于最新标记价格)
            "openOrderInitialMargin": "0.00000000", // 当前挂单所需起始保证金(基于最新标记价格)
            "crossWalletBalance": "23.72469206",  //全仓账户余额
            "crossUnPnl": "0.00000000" // 全仓持仓未实现盈亏
            "availableBalance": "23.72469206",       // 可用余额
            "maxWithdrawAmount": "23.72469206",     // 最大可转出余额
            "marginAvailable": true,   // 是否可用作联合保证金
            "updateTime": 1625474304765  //更新时间
        },
    ],
    "positions": [  // 头寸，将返回所有市场symbol。
        //根据用户持仓模式展示持仓方向，即单向模式下只返回BOTH持仓情况，双向模式下只返回 LONG 和 SHORT 持仓情况
        {
            "symbol": "BTCUSDT",  // 交易对
            "initialMargin": "0",   // 当前所需起始保证金(基于最新标记价格)
            "maintMargin": "0", //维持保证金
            "unrealizedProfit": "0.00000000",  // 持仓未实现盈亏
            "positionInitialMargin": "0",  // 持仓所需起始保证金(基于最新标记价格)
            "openOrderInitialMargin": "0",  // 当前挂单所需起始保证金(基于最新标记价格)
            "leverage": "100",  // 杠杆倍率
            "isolated": true,  // 是否是逐仓模式
            "entryPrice": "0.00000",  // 持仓成本价
            "maxNotional": "250000",  // 当前杠杆下用户可用的最大名义价值
            "bidNotional": "0",  // 买单净值，忽略
          	"askNotional": "0",  // 买单净值，忽略
            "positionSide": "BOTH",  // 持仓方向
            "positionAmt": "0",      // 持仓数量
            "updateTime": 0         // 更新时间
        }
    ]
}
*/
func QueryAccount() (*futures.Account, error) {
	c := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := c.NewGetAccountService().Do(ctx)
	if err != nil {
		return nil, err
	}

	//log.Println(toJson(res))
	for _, asset := range res.Assets {
		if Str2Float64(asset.WalletBalance) != 0 {
			log.Println("account asset", toJson(asset))
		}
	}
	for _, position := range res.Positions {
		if Str2Float64(position.PositionAmt) != 0 {
			log.Println("account position", toJson(position))
		}

	}

	return res, nil
}

/*
QueryAccountBalance /fapi/v2/balance 账户余额

        "accountAlias": "SgsR",    // 账户唯一识别码
        "asset": "USDT",        // 资产
        "balance": "122607.35137903",   // 总余额
        "crossWalletBalance": "23.72469206", // 全仓余额
        "crossUnPnl": "0.00000000"  // 全仓持仓未实现盈亏
        "availableBalance": "23.72469206",       // 下单可用余额
        "maxWithdrawAmount": "23.72469206",     // 最大可转出余额
        "marginAvailable": true,    // 是否可用作联合保证金
        "updateTime": 1617939110373
*/
func QueryAccountBalance() {
	balances, err := NewClient().NewGetBalanceService().Do(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	for _, balance := range balances {
		if Str2Float64(balance.Balance) != 0 {
			log.Println("balance", toJson(balance))
		}
	}
	//log.Println(toJsonIndent(balances))
}

func QueryAccountPositions() ([]*futures.AccountPosition, error) {
	account, err := QueryAccount()
	if err != nil {
		return nil, err
	}

	var positions []*futures.AccountPosition
	for _, p := range account.Positions {
		if Str2Float64(p.PositionAmt) != 0 {
			positions = append(positions, p)
		}
	}
	return positions, nil
}
func QueryAccountAssets() ([]*futures.AccountAsset, error) {
	account, err := QueryAccount()
	if err != nil {
		return nil, err
	}

	var assets []*futures.AccountAsset
	for _, a := range account.Assets {
		if Str2Float64(a.WalletBalance) != 0 {
			assets = append(assets, a)
		}
	}
	return assets, nil
}
