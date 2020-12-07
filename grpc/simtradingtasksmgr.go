package tradingdb2grpc

import (
	"context"

	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
)

// FuncOnTaskEnd -
type FuncOnTaskEnd func(context.Context, *tradingpb.RequestSimTrading, *tradingpb.ReplySimTrading, error)

// SimTradingTask -
type SimTradingTask struct {
	Req       *tradingpb.RequestSimTrading
	OnTaskEnd FuncOnTaskEnd
}

// SimTradingTasksMgr -
type SimTradingTasksMgr struct {
	MapTasks map[int32]*SimTradingTask
}
