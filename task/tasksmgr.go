package tradingdb2task

import (
	"bytes"
	"sync"

	"github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// TasksMgr - TasksMgr
type TasksMgr struct {
	mapTasks   map[interface{}]*Task
	mutex      sync.Mutex
	lstKeys    []interface{}
	lstRunning []interface{}
}

// NewTasksMgr - new TasksMgr
func NewTasksMgr() *TasksMgr {
	return &TasksMgr{
		mapTasks: make(map[interface{}]*Task),
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
	_, isok := mgr.mapTasks[buf]
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
		_, isok := mgr.mapTasks[buf]
		if !isok {
			mgr.mapTasks[buf] = &Task{
				Params: params,
			}

			mgr.lstKeys = append(mgr.lstKeys, buf)
		}

		return nil
	}

	task, isok := mgr.mapTasks[buf]
	if isok {
		task.lstFunc = append(task.lstFunc, onEnd)
	} else {
		mgr.mapTasks[buf] = &Task{
			Params:  params,
			lstFunc: []FuncOnTaskEnd{onEnd},
		}

		mgr.lstKeys = append(mgr.lstKeys, buf)
	}

	return nil
}

// func (mgr *TasksMgr) delKey(buf []byte) {
// 	for i, v := range mgr.lstKeys {
// 		bv, isok := v.([]byte)
// 		if isok {
// 			if bytes.Equal(bv, buf) {
// 				mgr.lstKeys = append(mgr.lstKeys[:i], mgr.lstKeys[i+1:]...)

// 				return
// 			}
// 		}
// 	}
// }

func (mgr *TasksMgr) delRunning(buf []byte) {
	for i, v := range mgr.lstRunning {
		bv, isok := v.([]byte)
		if isok {
			if bytes.Equal(bv, buf) {
				mgr.lstRunning = append(mgr.lstRunning[:i], mgr.lstRunning[i+1:]...)

				return
			}
		}
	}
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

	task, isok := mgr.mapTasks[buf]
	if isok {
		for _, v := range task.lstFunc {
			v(task)
		}

		delete(mgr.mapTasks, buf)
		mgr.delRunning(buf)
	}

	return nil
}

func (mgr *TasksMgr) StartTask(onStart FuncOnTaskStart) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	for len(mgr.lstKeys) > 0 {
		buf := mgr.lstKeys[0]
		mgr.lstKeys = mgr.lstKeys[1:]

		task, isok := mgr.mapTasks[buf]
		if isok {
			mgr.lstRunning = append(mgr.lstRunning, buf)
			onStart(task)

			return nil
		}
	}

	return nil
}
