package tradingdb2task

import (
	"sync"

	"github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// TasksMgr - TasksMgr
type TasksMgr struct {
	MapTasks map[interface{}]*Task
	mutex    sync.Mutex
}

// NewTasksMgr - new TasksMgr
func NewTasksMgr() *TasksMgr {
	return &TasksMgr{
		MapTasks: make(map[interface{}]*Task),
	}
}

func (mgr *TasksMgr) HasTask(params *tradingpb.SimTradingParams) bool {
	buf, err := proto.Marshal(params)
	if err != nil {
		tradingdb2utils.Warn("TasksMgr.HasTask:Marshal",
			zap.Error(err))

		return false
	}

	mgr.mutex.Lock()
	_, isok := mgr.MapTasks[buf]
	mgr.mutex.Unlock()

	return isok
}

func (mgr *TasksMgr) AddTask(params *tradingpb.SimTradingParams, onEnd FuncOnTaskEnd) error {
	buf, err := proto.Marshal(params)
	if err != nil {
		tradingdb2utils.Warn("TasksMgr.AddTask:Marshal",
			zap.Error(err))

		return err
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	if onEnd == nil {
		_, isok := mgr.MapTasks[buf]
		if !isok {
			mgr.MapTasks[buf] = &Task{
				Params: params,
			}
		}

		return nil
	}

	task, isok := mgr.MapTasks[buf]
	if isok {
		task.lstFunc = append(task.lstFunc, onEnd)
	} else {
		mgr.MapTasks[buf] = &Task{
			Params:  params,
			lstFunc: []FuncOnTaskEnd{onEnd},
		}
	}

	return nil
}

func (mgr *TasksMgr) OnTaskEnd(params *tradingpb.SimTradingParams) error {
	buf, err := proto.Marshal(params)
	if err != nil {
		tradingdb2utils.Warn("TasksMgr.OnTaskEnd:Marshal",
			zap.Error(err))

		return err
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	task, isok := mgr.MapTasks[buf]
	if isok {
		for _, v := range task.lstFunc {
			v(task)
		}

		delete(mgr.MapTasks, buf)
	}

	return nil
}
