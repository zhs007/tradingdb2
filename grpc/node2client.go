package tradingdb2grpc

import (
	"context"

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
				logger.Error("Client.GetServerInfo:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			} else {
				tradingdb2utils.Error("Client.GetServerInfo:grpc.Dial",
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
			logger.Error("Client.GetServerInfo:Client.GetServerInfo",
				zap.Error(err))
		} else {
			tradingdb2utils.Error("Client.GetServerInfo:Client.GetServerInfo",
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

	return reply, nil
}
