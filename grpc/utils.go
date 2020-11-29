package tradingdb2grpc

import (
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
)

// GenCalcBaseline - generator SimTradingParams for baseline
func GenCalcBaseline(asset *tradingpb.Asset, startTs int64, endTs int64) *tradingpb.SimTradingParams {
	pmbuy := &tradingpb.BuyParams{
		PerInitMoney: 1,
	}

	pminit := &tradingpb.InitParams{
		Money: 10000,
	}

	ccbuy := &tradingpb.CtrlCondition{
		Indicator: "buyandhold",
	}

	strategy := &tradingpb.Strategy{
		Name:       "bah",
		Asset:      asset,
		Buy:        []*tradingpb.CtrlCondition{ccbuy},
		ParamsBuy:  pmbuy,
		ParamsInit: pminit,
	}

	params := &tradingpb.SimTradingParams{
		Assets:     []*tradingpb.Asset{asset},
		StartTs:    startTs,
		EndTs:      endTs,
		Strategies: []*tradingpb.Strategy{strategy},
	}

	return params
}
