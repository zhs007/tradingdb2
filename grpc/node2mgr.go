package tradingdb2grpc

import (
	"context"

	tradingdb2 "github.com/zhs007/tradingdb2"
	"github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
)

// Node2Mgr - Node2 manager
type Node2Mgr struct {
	Cfg   *tradingdb2.Config
	Nodes []*Node2Client
}

// NewNode2Mgr - new a Node2Mgr
func NewNode2Mgr(cfg *tradingdb2.Config) (*Node2Mgr, error) {
	mgr := &Node2Mgr{
		Cfg: cfg,
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
	for {
		cn = mgr.getFreeClient()
		if cn != nil {
			break
		}
	}

	return cn.CalcPNL(ctx, params, logger)
}
