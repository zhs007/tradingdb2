package tradingdb2grpc

import (
	"context"
	"time"

	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Node2Client - TradingNode2Client
type Node2Client struct {
	servAddr     string
	token        string
	conn         *grpc.ClientConn
	client       tradingpb.TradingNode2Client
	lastTaskNums int
	lastTs       int64
}

// NewNode2Client - new TradingNode2Client
func NewNode2Client(servAddr string, token string) (*Node2Client, error) {
	client := &Node2Client{
		servAddr:     servAddr,
		token:        token,
		lastTaskNums: 1,
	}

	return client, nil
}

// isFree - isFree
func (client *Node2Client) isFree() bool {
	return client.lastTaskNums > 0
}

// reset - reset
func (client *Node2Client) reset() {
	if client.conn != nil {
		client.conn.Close()
	}

	client.conn = nil
	client.client = nil
}

// GetServerInfo - getServerInfo
func (client *Node2Client) GetServerInfo(ctx context.Context, logger *zap.Logger) (*tradingpb.ReplyServerInfo, error) {

	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			if logger != nil {
				logger.Error("Node2Client.GetServerInfo:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			} else {
				tradingdb2utils.Error("Node2Client.GetServerInfo:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			}

			return nil, err
		}

		client.conn = conn
		client.client = tradingpb.NewTradingNode2Client(conn)
	}

	req := &tradingpb.RequestServerInfo{
		BasicRequest: &tradingpb.BasicRequestData{
			Token: client.token,
		},
	}

	reply, err := client.client.GetServerInfo(ctx, req)
	if err != nil {
		if logger != nil {
			logger.Error("Node2Client.GetServerInfo:Client.GetServerInfo",
				zap.Error(err))
		} else {
			tradingdb2utils.Error("Node2Client.GetServerInfo:Client.GetServerInfo",
				zap.Error(err))
		}

		// if error, reset
		client.reset()

		return nil, err
	}

	return reply, nil
}

// CalcPNL - calcPNL
func (client *Node2Client) CalcPNL(ctx context.Context, params *tradingpb.SimTradingParams, logger *zap.Logger) (*tradingpb.ReplyCalcPNL, error) {
	if client.lastTaskNums <= 0 {
		if time.Now().Unix()-client.lastTs < TradingNode2RequestOffTime {
			tradingdb2utils.Error("Node2Client.CalcPNL",
				zap.Error(ErrNodeNotFree))

			return nil, ErrNodeNotFree
		}

		client.lastTaskNums = 1
	}

	client.lastTaskNums--
	client.lastTs = time.Now().Unix()

	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			if logger != nil {
				logger.Error("Client.CalcPNL:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			} else {
				tradingdb2utils.Error("Client.CalcPNL:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			}

			return nil, err
		}

		client.conn = conn
		client.client = tradingpb.NewTradingNode2Client(conn)
	}

	req := &tradingpb.RequestCalcPNL{
		BasicRequest: &tradingpb.BasicRequestData{
			Token: client.token,
		},
		Params: params,
	}

	reply, err := client.client.CalcPNL(ctx, req)
	if err != nil {
		if logger != nil {
			logger.Error("Client.CalcPNL:Client.CalcPNL",
				zap.Error(err))
		} else {
			tradingdb2utils.Error("Client.CalcPNL:Client.CalcPNL",
				zap.Error(err))
		}

		// if error, reset
		client.reset()

		return nil, err
	}

	if reply.NodeInfo != nil {
		client.lastTaskNums = int(reply.NodeInfo.MaxTasks - reply.NodeInfo.CurTasks)
	} else {
		client.lastTaskNums++
	}

	return reply, nil
}
