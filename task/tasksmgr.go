package tradingdb2task

import (
	"bytes"
	"sync"
	"time"

	"github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// TasksMgr - TasksMgr
type TasksMgr struct {
	mapTasks          map[interface{}]*Task
	mutex             sync.Mutex
	lstKeys           []interface{}
	lstRunning        []interface{}
	latestTaskGroupID int
	mapTaskGroup      map[int]*TaskGroup
}

// NewTasksMgr - new TasksMgr
func NewTasksMgr() *TasksMgr {
	return &TasksMgr{
		mapTasks:          make(map[interface{}]*Task),
		latestTaskGroupID: 0,
		mapTaskGroup:      make(map[int]*TaskGroup),
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

func (mgr *TasksMgr) AddTask(taskGroupID int, params *tradingpb.SimTradingParams, onEnd FuncOnTaskEnd) error {
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
				Params:      params,
				TaskGroupID: taskGroupID,
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
			Params:      params,
			lstFunc:     []FuncOnTaskEnd{onEnd},
			TaskGroupID: taskGroupID,
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

func (mgr *TasksMgr) OnTaskEnd(result *tradingpb.TradingTaskResult) error {
	// buf, err := proto.Marshal(params)
	// if err != nil {
	// 	tradingdb2utils.Warn("TasksMgr.OnTaskEnd:Marshal",
	// 		zap.Error(err))

	// 	return err
	// }

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	task, isok := mgr.mapTasks[result.Task]
	if isok {
		task.PNL = result.Pnl

		for _, v := range task.lstFunc {
			v(task)
		}

		delete(mgr.mapTasks, result.Task)
		mgr.delRunning(result.Task)

		return nil
	}

	tradingdb2utils.Warn("TasksMgr.OnTaskEnd",
		zap.Error(ErrInvalidTask))

	return ErrInvalidTask
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

			task.StartTs = time.Now().Unix()

			onStart(task)

			return nil
		}
	}

	onStart(nil)

	return nil
}

func (mgr *TasksMgr) NewTaskGroup() int {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	mgr.latestTaskGroupID++

	mgr.mapTaskGroup[mgr.latestTaskGroupID] = &TaskGroup{
		TaskGroupID: mgr.latestTaskGroupID,
		StartTs:     time.Now().Unix(),
	}

	return mgr.latestTaskGroupID
}

func (mgr *TasksMgr) IsTaskGroupFinished(taskGroupID int) bool {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	for _, v := range mgr.mapTasks {
		if v.TaskGroupID == taskGroupID {
			return false
		}
	}

	return true
}

func (mgr *TasksMgr) LogTaskGroup(taskGroupID int, str string) {
	logts := time.Now().Unix()

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	tg, isok := mgr.mapTaskGroup[taskGroupID]
	if isok {
		tradingdb2utils.Info(str,
			zap.Int("taskGroupID", taskGroupID),
			zap.Duration("Running", time.Duration(logts-tg.StartTs)))
	}
}

func (mgr *TasksMgr) WaitTaskGroupFinished(taskGroupID int) {
	logts := time.Now().Unix()

	mgr.LogTaskGroup(taskGroupID, "TasksMgr.WaitTaskGroupFinished")

	for {
		if mgr.IsTaskGroupFinished(taskGroupID) {
			break
		}

		time.Sleep(5 * time.Second)

		ts := time.Now().Unix()
		if ts-logts >= int64(30*time.Second) {
			mgr.LogTaskGroup(taskGroupID, "TasksMgr.WaitTaskGroupFinished...")

			logts = ts
		}
	}
}
