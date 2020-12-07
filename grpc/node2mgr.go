package tradingdb2grpc

import (
	"context"
	"sync"
	"time"

	tradingdb2 "github.com/zhs007/tradingdb2"
	"github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
)

// FuncOnCalcPNLEnd -
type FuncOnCalcPNLEnd func(int, *tradingpb.SimTradingParams, *tradingpb.ReplyCalcPNL, error)

// Node2Task -
type Node2Task struct {
	TaskIndex int
	Params    *tradingpb.SimTradingParams
	OnEnd     FuncOnCalcPNLEnd
}

// Node2TaskResult -
type Node2TaskResult struct {
	Task  *Node2Task
	Reply *tradingpb.ReplyCalcPNL
	Err   error
}

// Node2Mgr - Node2 manager
type Node2Mgr struct {
	Cfg            *tradingdb2.Config
	Nodes          []*Node2Client
	tasks          []*Node2Task
	tasksRunning   map[*Node2Task]*Node2Task
	mutexTasks     sync.Mutex
	chanTaskResult chan *Node2TaskResult
	chanRelease    chan int
	chanAddTask    chan int
	isRunning      bool
}

// NewNode2Mgr - new a Node2Mgr
func NewNode2Mgr(cfg *tradingdb2.Config) (*Node2Mgr, error) {
	mgr := &Node2Mgr{
		Cfg:            cfg,
		tasksRunning:   make(map[*Node2Task]*Node2Task),
		chanTaskResult: make(chan *Node2TaskResult),
		chanRelease:    make(chan int),
		chanAddTask:    make(chan int),
	}

	for _, v := range cfg.Nodes {
		cc, err := NewNode2Client(v.Host, v.Token)
		if err != nil {
			tradingdb2utils.Error("NewServ.NewNode2Mgr",
				zap.Error(err))

			return nil, err
		}

		mgr.Nodes = append(mgr.Nodes, cc)
	}

	return mgr, nil
}

// getFreeClient - get a free client
func (mgr *Node2Mgr) getFreeClient() *Node2Client {
	for _, v := range mgr.Nodes {
		if v.isFree() {
			return v
		}
	}

	return nil
}

// CalcPNL - calcPNL
func (mgr *Node2Mgr) CalcPNL(ctx context.Context, params *tradingpb.SimTradingParams, logger *zap.Logger) (*tradingpb.ReplyCalcPNL, error) {
	var cn *Node2Client
	ts := time.Now().Unix()
	for {
		cn = mgr.getFreeClient()
		if cn != nil {
			break
		}

		time.Sleep(time.Second)

		curts := time.Now().Unix()
		if curts-ts > WaitTradingNode2Time {
			break
		}
	}

	if cn != nil {
		return cn.CalcPNL(ctx, params, logger)
	}

	tradingdb2utils.Error("Node2Mgr.CalcPNL",
		zap.Error(ErrNodeNotFree))

	return nil, ErrNodeNotFree
}

// CalcPNL2 - calcPNL
func (mgr *Node2Mgr) CalcPNL2(ctx context.Context, params *tradingpb.SimTradingParams, logger *zap.Logger) (*tradingpb.ReplyCalcPNL, error) {
	chanCur := make(chan *Node2TaskResult)

	mgr.AddTask(-1, params, func(taskIndex int, params *tradingpb.SimTradingParams, reply *tradingpb.ReplyCalcPNL, err error) {
		chanCur <- &Node2TaskResult{
			Reply: reply,
			Err:   err,
		}
	}, logger)

	ret := <-chanCur

	return ret.Reply, ret.Err
}

// AddTask - add a task
func (mgr *Node2Mgr) AddTask(taskIndex int, params *tradingpb.SimTradingParams,
	onEnd FuncOnCalcPNLEnd, logger *zap.Logger) error {

	task := &Node2Task{
		TaskIndex: taskIndex,
		Params:    params,
		OnEnd:     onEnd,
	}

	mgr.mutexTasks.Lock()
	defer mgr.mutexTasks.Unlock()

	mgr.tasks = append(mgr.tasks, task)

	mgr.chanAddTask <- 0

	return nil
}

// AddTask - add a task
func (mgr *Node2Mgr) onMain() error {
	mgr.isRunning = true

	isStop := false

	for {
		select {
		case <-mgr.chanAddTask:
			if !isStop {
				mgr.nextTask(context.Background())
			}

		case result := <-mgr.chanTaskResult:
			mgr.onTaskEnd(result)

			if !isStop {
				mgr.nextTask(context.Background())
			}

		case <-mgr.chanRelease:
			isStop = true

		}

		if isStop && !mgr.hasRunningTasks() {
			mgr.isRunning = false

			return nil
		}
	}
}

// Start -
func (mgr *Node2Mgr) Start() error {
	go func() {
		mgr.onMain()
	}()

	return nil
}

// Stop -
func (mgr *Node2Mgr) Stop() error {
	mgr.chanRelease <- 0

	for {
		if !mgr.isRunning {
			break
		}

		time.Sleep(time.Second)
	}

	return nil
}

// onTaskEnd -
func (mgr *Node2Mgr) onTaskEnd(result *Node2TaskResult) error {
	mgr.mutexTasks.Lock()
	defer mgr.mutexTasks.Unlock()

	if result != nil {
		if result.Task != nil {
			if result.Task.OnEnd != nil {
				result.Task.OnEnd(result.Task.TaskIndex, result.Task.Params, result.Reply, result.Err)
			}

			_, isok := mgr.tasksRunning[result.Task]
			if isok {
				delete(mgr.tasksRunning, result.Task)

				return nil
			}

			tradingdb2utils.Error("Node2Mgr.onTaskEnd",
				zap.Error(ErrNoRuningTaskInNode2Mgr))
		}
	}

	return ErrNoRuningTaskInNode2Mgr
}

// nextTask -
func (mgr *Node2Mgr) nextTask(ctx context.Context) error {
	mgr.mutexTasks.Lock()
	defer mgr.mutexTasks.Unlock()

	client := mgr.getFreeClient()
	if client == nil {
		return nil
	}

	if len(mgr.tasks) > 0 {
		curtask := mgr.tasks[0]

		_, isok := mgr.tasksRunning[curtask]
		if isok {
			tradingdb2utils.Error("Node2Mgr.nextTask:tasksRunning",
				zap.Error(ErrDuplicateRuningTaskInNode2Mgr))

			return ErrDuplicateRuningTaskInNode2Mgr
		}

		mgr.tasks = mgr.tasks[1:]
		mgr.tasksRunning[curtask] = curtask

		go func() {
			mgr.runTask(ctx, client, curtask)
		}()
	}

	return nil
}

// runTask -
func (mgr *Node2Mgr) runTask(ctx context.Context, client *Node2Client, task *Node2Task) {
	if client != nil {
		reply, err := client.CalcPNL(ctx, task.Params, nil)
		if err != nil {
			tradingdb2utils.Error("Node2Mgr.runTask:CalcPNL",
				zap.Error(err))

			mgr.chanTaskResult <- &Node2TaskResult{
				Task: task,
				Err:  err,
			}
		} else {
			mgr.chanTaskResult <- &Node2TaskResult{
				Task:  task,
				Reply: reply,
			}
		}

		return
	}

	tradingdb2utils.Error("Node2Mgr.runTask",
		zap.Error(ErrNoNode))
}

// hasRunningTasks -
func (mgr *Node2Mgr) hasRunningTasks() bool {
	mgr.mutexTasks.Lock()
	defer mgr.mutexTasks.Unlock()

	return len(mgr.tasksRunning) > 0
}
