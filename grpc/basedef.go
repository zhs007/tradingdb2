package tradingdb2grpc

// TradingNode2TimeOut -
const TradingNode2TimeOut = 60

// WaitTradingNode2Time -
const WaitTradingNode2Time = 10

// ReqTradingTask3TimeOut -
const ReqTradingTask3TimeOut = 10 * 60

type reqTasks3Result struct {
	err      error
	isEnd    bool
	taskNums int
}
