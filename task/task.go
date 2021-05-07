package tradingdb2task

import "github.com/zhs007/tradingdb2/tradingpb"

// FuncOnTaskEnd - used in BatchCandles
type FuncOnTaskEnd func(task *Task) error

// Task - Task
type Task struct {
	Params  *tradingpb.SimTradingParams
	lstFunc []FuncOnTaskEnd
}
