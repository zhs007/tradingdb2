package tradingdb2task

import (
	"encoding/hex"
	"sync"
	"time"

	"github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// TasksMgr - TasksMgr
type TasksMgr struct {
	mapTasks          map[string]*Task
	mutex             sync.Mutex
	lstKeys           []string
	lstRunning        []string
	latestTaskGroupID int
	mapTaskGroup      map[int]*TaskGroup
	lstHistory        []*TaskGroup
}

// NewTasksMgr - new TasksMgr
func NewTasksMgr() *TasksMgr {
	return &TasksMgr{
		mapTasks:          make(map[string]*Task),
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
	_, isok := mgr.mapTasks[hex.EncodeToString(buf)]
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

	key := hex.EncodeToString(buf)

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	tg, isok := mgr.mapTaskGroup[taskGroupID]
	if isok {
		tg.MaxTaskNums++
	}

	if onEnd == nil {
		_, isok := mgr.mapTasks[key]
		if !isok {
			mgr.mapTasks[key] = &Task{
				Params:      params,
				TaskGroupID: taskGroupID,
			}

			mgr.lstKeys = append(mgr.lstKeys, key)
		}

		return nil
	}

	task, isok := mgr.mapTasks[key]
	if isok {
		task.lstFunc = append(task.lstFunc, onEnd)
	} else {
		mgr.mapTasks[key] = &Task{
			Params:      params,
			lstFunc:     []FuncOnTaskEnd{onEnd},
			TaskGroupID: taskGroupID,
		}

		mgr.lstKeys = append(mgr.lstKeys, key)
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

func (mgr *TasksMgr) delRunning(key string) {
	// key := hex.EncodeToString(buf)

	for i, v := range mgr.lstRunning {
		// bv, isok := v.([]byte)
		// if isok {
		if v == key {
			mgr.lstRunning = append(mgr.lstRunning[:i], mgr.lstRunning[i+1:]...)

			return
		}
		// }
	}

	tradingdb2utils.Warn("TasksMgr.delRunning",
		zap.String("key", key),
		zap.Error(ErrInvalidTaskKey))
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

	key := hex.EncodeToString(result.Task)

	if result.Err != "" {
		tradingdb2utils.Warn("TasksMgr.OnTaskEnd:Result",
			zap.String("Error", result.Err))

		// 重新加回队列
		mgr.delRunning(key)
		mgr.lstKeys = append(mgr.lstKeys, key)

		return ErrTaskFail
	}

	task, isok := mgr.mapTasks[key]
	if isok {
		task.PNL = result.Pnl

		for _, v := range task.lstFunc {
			v(task)
		}

		delete(mgr.mapTasks, key)
		mgr.delRunning(key)

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
		key := mgr.lstKeys[0]
		mgr.lstKeys = mgr.lstKeys[1:]

		// key := hex.EncodeToString(buf)

		task, isok := mgr.mapTasks[key]
		if isok {
			mgr.lstRunning = append(mgr.lstRunning, key)

			task.StartTs = time.Now().Unix()

			onStart(task)

			return nil
		} else {
			tradingdb2utils.Warn("TasksMgr.StartTask",
				zap.String("key", key),
				zap.Error(ErrInvalidTaskKey))
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
		TaskGroupID:  mgr.latestTaskGroupID,
		StartTs:      time.Now().Unix(),
		LastTaskNums: -1,
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

	tg, isok := mgr.mapTaskGroup[taskGroupID]
	if isok {
		mgr.addHistory(tg)

		delete(mgr.mapTaskGroup, taskGroupID)
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
			zap.Duration("Running", time.Duration((logts-tg.StartTs)*int64(time.Second))))
	}
}

func (mgr *TasksMgr) WaitTaskGroupFinished(taskGroupID int) {
	logts := time.Now().Unix()

	tg, isok := mgr.mapTaskGroup[taskGroupID]
	if isok {
		tg.IsRecvEnd = true
	}

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

func (mgr *TasksMgr) GetTaskGroups() []TaskGroup {
	arr := []TaskGroup{}

	curts := time.Now().Unix()

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	for _, tg := range mgr.mapTaskGroup {
		tg.RunningTime = curts - tg.StartTs
		if tg.RunningTime <= 0 {
			tg.RunningTime = 0
			tg.StartTs = curts
		}

		tg.LastTaskNums = 0
		tg.RunningTaskNums = 0
		for _, t := range mgr.mapTasks {
			if t.TaskGroupID == tg.TaskGroupID {
				tg.LastTaskNums++

				if mgr.isInRunning(t.Params) {
					tg.RunningTaskNums++
				}
			}
		}

		if tg.MaxTaskNums <= tg.LastTaskNums {
			tg.MaxTaskNums = tg.LastTaskNums
		}

		tg.RunningTimeStr = time.Duration(tg.RunningTime * int64(time.Second)).String()

		if tg.MaxTaskNums == tg.LastTaskNums || tg.RunningTime <= 0 {
			tg.LastTime = -1
			tg.LastTimeStr = "Please wait a moment."
		} else {
			tg.LastTime = int64(tg.LastTaskNums) * tg.RunningTime / int64(tg.MaxTaskNums-tg.LastTaskNums)
			tg.LastTimeStr = time.Duration(tg.LastTime * int64(time.Second)).String()
		}

		if tg.LastTaskNums > 0 {
			arr = append(arr, *tg)
		}
	}

	return arr
}

func (mgr *TasksMgr) addHistory(tg *TaskGroup) {
	curts := time.Now().Unix()

	tg.RunningTime = curts - tg.StartTs
	tg.RunningTimeStr = time.Duration(tg.RunningTime * int64(time.Second)).String()
	tg.LastTime = 0
	tg.LastTimeStr = ""
	tg.LastTaskNums = 0
	tg.RunningTaskNums = 0

	mgr.lstHistory = append(mgr.lstHistory, tg)
}

func (mgr *TasksMgr) RecvHistory() []TaskGroup {
	if len(mgr.lstHistory) == 0 {
		return nil
	}

	arr := []TaskGroup{}

	for _, v := range mgr.lstHistory {
		if v.MaxTaskNums > 0 {
			arr = append(arr, *v)
		}
	}

	mgr.lstHistory = nil

	return arr
}

func (mgr *TasksMgr) GetLastTasks(taskGroupID int) []ShowTaskObj {
	arr := []ShowTaskObj{}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	for k, v := range mgr.mapTasks {
		if v.TaskGroupID == taskGroupID {
			arr = append(arr, ShowTaskObj{
				Params: v.Params,
				Key:    k,
			})
		}
	}

	return arr
}

func (mgr *TasksMgr) ResetTaskKeyList(lst *TaskKeyList) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	for _, v := range lst.lst {
		mgr.resetTaskKey(v)
	}
}

func (mgr *TasksMgr) resetTaskKey(key string) error {
	_, isok := mgr.mapTasks[key]
	if !isok {
		tradingdb2utils.Warn("TasksMgr:resetTaskKey",
			zap.String("key", key),
			zap.Error(ErrNoKey))

		return ErrNoKey
	}

	for i, v := range mgr.lstRunning {
		if v == key {
			mgr.lstRunning = append(mgr.lstRunning[:i], mgr.lstRunning[i+1:]...)

			if mgr.addKey(key) {
				return nil
			}

			tradingdb2utils.Warn("TasksMgr:resetTaskKey:addKey",
				zap.String("key", key),
				zap.Error(ErrDuplicateKey))

			return ErrDuplicateKey
		}
	}

	if mgr.addKey(key) {
		return nil
	}

	tradingdb2utils.Warn("TasksMgr:resetTaskKey:addKey:non-running",
		zap.String("key", key),
		zap.Error(ErrDuplicateKey))

	return ErrDuplicateKey
}

func (mgr *TasksMgr) addKey(key string) bool {
	for _, v := range mgr.lstKeys {
		if v == key {
			return false
		}
	}

	mgr.lstKeys = append(mgr.lstKeys, key)

	return true
}

func (mgr *TasksMgr) isInRunning(params *tradingpb.SimTradingParams) bool {
	buf, err := proto.Marshal(params)
	if err != nil {
		tradingdb2utils.Warn("TasksMgr.isInRunning:Marshal",
			zap.Error(err))

		return false
	}

	key := hex.EncodeToString(buf)

	for _, v := range mgr.lstRunning {
		if v == key {
			return true
		}
	}

	return false
}
