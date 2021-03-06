package tradingdb2task

import "github.com/zhs007/tradingdb2/tradingpb"

// FuncOnTaskEnd - on task end
type FuncOnTaskEnd func(task *Task) error

// FuncOnTaskStart - on task start
type FuncOnTaskStart func(task *Task) error

// Task - Task
type Task struct {
	Params      *tradingpb.SimTradingParams
	PNL         *tradingpb.PNLData
	lstFunc     []FuncOnTaskEnd
	StartTs     int64
	TaskGroupID int
}

// ShowTaskObj - ShowTaskObj
type ShowTaskObj struct {
	Params *tradingpb.SimTradingParams `json:"params"`
	Key    string                      `json:"key"`
}
