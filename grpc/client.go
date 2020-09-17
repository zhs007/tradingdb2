package tradingdb2grpc

import (
	"context"

	tradingdb2pb "github.com/zhs007/tradingdb2/pb"
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
func (client *Client) UpdCandles(ctx context.Context, candles *tradingdb2pb.Candles) (
	*tradingdb2pb.ReplyUpdCandles, error) {

	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		client.conn = conn
		client.client = tradingdb2pb.NewTradingDB2ServiceClient(conn)
	}

	reply, err := client.client.UpdCandles(ctx, &tradingdb2pb.RequestUpdCandles{
		Token:   client.token,
		Candles: candles,
	})
	if err != nil {
		// if error, reset
		client.reset()

		return nil, err
	}

	return reply, nil
}

// GetCandles - get candles
func (client *Client) GetCandles(ctx context.Context, market string, symbol string, tag string) (
	*tradingdb2pb.ReplyGetCandles, error) {

	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		client.conn = conn
		client.client = tradingdb2pb.NewTradingDB2ServiceClient(conn)
	}

	reply, err := client.client.GetCandles(ctx, &tradingdb2pb.RequestGetCandles{
		Token:  client.token,
		Market: market,
		Symbol: symbol,
		Tag:    tag,
	})
	if err != nil {
		// if error, reset
		client.reset()

		return nil, err
	}

	return reply, nil
}
