package tradingdb2grpc

import (
	"sync"
	"time"

	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
)

// FuncOnSimTradingTaskEnd -
type FuncOnSimTradingTaskEnd func(*tradingpb.RequestSimTrading, *tradingpb.ReplySimTrading, error)

// SimTradingTask -
type SimTradingTask struct {
	Req   *tradingpb.RequestSimTrading
	OnEnd FuncOnSimTradingTaskEnd
}

// SimTradingTasksMgr -
type SimTradingTasksMgr struct {
	mapTasks    map[int32]*SimTradingTask
	mutexTasks  sync.Mutex
	chanResult  chan *SimTradingTask
	chanRelease chan int
	isRunning   bool
}

// NewSimTradingTasksMgr - new a SimTradingTasksMgr
func NewSimTradingTasksMgr() *SimTradingTasksMgr {
	return &SimTradingTasksMgr{
		mapTasks:    make(map[int32]*SimTradingTask),
		chanResult:  make(chan *SimTradingTask),
		chanRelease: make(chan int),
	}
}

// AddTask -
func (mgr *SimTradingTasksMgr) AddTask(mgrNode *Node2Mgr, req *tradingpb.RequestSimTrading, onEnd FuncOnSimTradingTaskEnd) error {
	mgr.mutexTasks.Lock()

	_, isok := mgr.mapTasks[req.Index]
	if isok {
		mgr.mutexTasks.Unlock()

		tradingdb2utils.Error("SimTradingTasksMgr.AddTask",
			zap.Error(ErrDuplicateTaskIndex))

		return ErrDuplicateTaskIndex
	}

	curtask := &SimTradingTask{
		Req:   req,
		OnEnd: onEnd,
	}

	mgr.mapTasks[req.Index] = curtask

	defer mgr.mutexTasks.Unlock()

	mgrNode.AddTask(int(req.Index), req.Params, func(taskIndex int, params *tradingpb.SimTradingParams, reply *tradingpb.ReplyCalcPNL, err error) {
		if reply != nil {
			onEnd(req, &tradingpb.ReplySimTrading{
				Pnl: reply.Pnl,
			}, err)
		} else {
			onEnd(req, nil, err)
		}

		mgr.chanResult <- curtask
	}, nil)

	return nil
}

// IsRunning -
func (mgr *SimTradingTasksMgr) IsRunning() bool {
	return mgr.isRunning
}

// isFinished -
func (mgr *SimTradingTasksMgr) isFinished() bool {
	mgr.mutexTasks.Lock()
	defer mgr.mutexTasks.Unlock()

	return len(mgr.mapTasks) <= 0
}

// Start -
func (mgr *SimTradingTasksMgr) Start() {
	go func() {
		mgr.onMain()
	}()
}

// Stop - 会等待队列任务全部完成
func (mgr *SimTradingTasksMgr) Stop() {
	mgr.chanRelease <- 0

	for {
		if !mgr.isRunning {
			break
		}

		time.Sleep(time.Second)
	}
}

// onMain -
func (mgr *SimTradingTasksMgr) onMain() {
	mgr.isRunning = true

	isStop := false

	for {
		select {
		case task := <-mgr.chanResult:
			mgr.onTaskEnd(task)

		case <-mgr.chanRelease:
			isStop = true
		}

		if isStop && mgr.isFinished() {
			mgr.isRunning = false

			return
		}
	}
}

// onTaskEnd -
func (mgr *SimTradingTasksMgr) onTaskEnd(curtask *SimTradingTask) {
	if curtask != nil {
		mgr.mutexTasks.Lock()
		defer mgr.mutexTasks.Unlock()

		_, isok := mgr.mapTasks[curtask.Req.Index]
		if !isok {
			tradingdb2utils.Error("SimTradingTasksMgr.AddTask",
				zap.Error(ErrDuplicateTaskIndex))

			return
		}

		delete(mgr.mapTasks, curtask.Req.Index)
	}
}
