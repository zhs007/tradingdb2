package tradingdb2grpc

import (
	"context"

	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"

	"google.golang.org/grpc/peer"
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

func setIgnoreReplySimTrading(reply *tradingpb.ReplySimTrading) {
	for _, v := range reply.Pnl {
		v.Total.Values = nil
	}
}

func GetPeerAddr(ctx context.Context) string {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		tradingdb2utils.Warn("GetPeerAddr:FromContext")

		return ""
	}

	return pr.Addr.String()
}
