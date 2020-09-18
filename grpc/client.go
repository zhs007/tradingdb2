package tradingdb2grpc

import (
	"context"
	"io"

	tradingdb2 "github.com/zhs007/tradingdb2"
	tradingdb2pb "github.com/zhs007/tradingdb2/pb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Client - TradingDB2ServiceClient
type Client struct {
	servAddr string
	token    string
	conn     *grpc.ClientConn
	client   tradingdb2pb.TradingDB2ServiceClient
}

// NewClient - new TradingDB2ServiceClient
func NewClient(servAddr string, token string) (*Client, error) {
	client := &Client{
		servAddr: servAddr,
		token:    token,
	}

	return client, nil
}

// reset - reset
func (client *Client) reset() {
	if client.conn != nil {
		client.conn.Close()
	}

	client.conn = nil
	client.client = nil
}

// UpdCandles - update candles
func (client *Client) UpdCandles(ctx context.Context, candles *tradingdb2pb.Candles, batchNums int) (
	*tradingdb2pb.ReplyUpdCandles, error) {

	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			tradingdb2utils.Error("Client.UpdCandles:grpc.Dial",
				zap.String("server address", client.servAddr),
				zap.Error(err))

			return nil, err
		}

		client.conn = conn
		client.client = tradingdb2pb.NewTradingDB2ServiceClient(conn)
	}

	stream, err := client.client.UpdCandles(ctx)
	if err != nil {
		tradingdb2utils.Error("Client.UpdCandles:Client.UpdCandles",
			zap.Error(err))

		// if error, reset
		client.reset()

		return nil, err
	}

	sentnums := 0
	err = tradingdb2.BatchCandles(candles, batchNums, func(lst *tradingdb2pb.Candles) error {
		cc := &tradingdb2pb.RequestUpdCandles{
			Token:   client.token,
			Candles: lst,
		}

		if sentnums == 0 {
			cc.Candles.Market = candles.Market
			cc.Candles.Symbol = candles.Symbol
			cc.Candles.Tag = candles.Tag
		}

		err := stream.Send(cc)
		if err != nil {
			return err
		}

		sentnums += len(lst.Candles)

		return nil
	})
	if err != nil {
		tradingdb2utils.Error("Client.UpdCandles:BatchCandles",
			zap.Int("sent nums", sentnums),
			zap.Error(err))

		// if error, reset
		client.reset()

		return nil, err
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		tradingdb2utils.Error("Client.UpdCandles:stream.CloseAndRecv",
			zap.Int("sent nums", sentnums),
			zap.Error(err))

		// if error, reset
		client.reset()

		return nil, err
	}

	return reply, nil
}

// GetCandles - get candles
func (client *Client) GetCandles(ctx context.Context, market string, symbol string, tag string) (
	*tradingdb2pb.Candles, error) {

	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			tradingdb2utils.Error("Client.UpdCandles:grpc.Dial",
				zap.String("server address", client.servAddr),
				zap.Error(err))

			return nil, err
		}

		client.conn = conn
		client.client = tradingdb2pb.NewTradingDB2ServiceClient(conn)
	}

	stream, err := client.client.GetCandles(ctx, &tradingdb2pb.RequestGetCandles{
		Token:  client.token,
		Market: market,
		Symbol: symbol,
		Tag:    tag,
	})
	if err != nil {
		tradingdb2utils.Error("Client.UpdCandles:client.GetCandles",
			zap.Error(err))

		// if error, reset
		client.reset()

		return nil, err
	}

	candles := &tradingdb2pb.Candles{
		Market: market,
		Symbol: symbol,
		Tag:    tag,
	}
	times := 0

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return candles, nil
			}

			tradingdb2utils.Error("Client.UpdCandles:stream.Recv",
				zap.Int("times", times),
				zap.Int("candle nums", len(candles.Candles)),
				zap.Error(err))

			// if error, reset
			client.reset()

			return nil, err
		}

		times++

		tradingdb2.MergeCandles(candles, req.Candles)
	}
}
