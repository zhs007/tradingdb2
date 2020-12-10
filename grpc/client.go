package tradingdb2grpc

import (
	"context"
	"io"

	tradingdb2 "github.com/zhs007/tradingdb2"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// FuncOnGetSymbols - onGetSymbols(*tradingpb.SymbolInfo)
type FuncOnGetSymbols func(*tradingpb.SymbolInfo)

// Client - TradingDB2Client
type Client struct {
	servAddr string
	token    string
	conn     *grpc.ClientConn
	client   tradingpb.TradingDB2Client
}

// NewClient - new TradingDB2Client
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
func (client *Client) UpdCandles(ctx context.Context, candles *tradingpb.Candles, batchNums int, logger *zap.Logger) (
	*tradingpb.ReplyUpdCandles, error) {

	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			if logger != nil {
				logger.Error("Client.UpdCandles:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			} else {
				tradingdb2utils.Error("Client.UpdCandles:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			}

			return nil, err
		}

		client.conn = conn
		client.client = tradingpb.NewTradingDB2Client(conn)
	}

	stream, err := client.client.UpdCandles(ctx)
	if err != nil {
		if logger != nil {
			logger.Error("Client.UpdCandles:Client.UpdCandles",
				zap.Error(err))
		} else {
			tradingdb2utils.Error("Client.UpdCandles:Client.UpdCandles",
				zap.Error(err))
		}

		// if error, reset
		client.reset()

		return nil, err
	}

	sentnums := 0
	err = tradingdb2.BatchCandles(candles, batchNums, func(lst *tradingpb.Candles) error {
		cc := &tradingpb.RequestUpdCandles{
			Candles: lst,
			BasicRequest: &tradingpb.BasicRequestData{
				Token: client.token,
			},
		}

		if sentnums == 0 {
			cc.Candles.Market = candles.Market
			cc.Candles.Symbol = candles.Symbol
			cc.Candles.Tag = candles.Tag
		}

		err := stream.Send(cc)
		if err != nil {
			if logger != nil {
				logger.Error("Client.UpdCandles:stream.Send",
					zap.Error(err))
			} else {
				tradingdb2utils.Error("Client.UpdCandles:stream.Send",
					zap.Error(err))
			}

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
		if logger != nil {
			logger.Error("Client.UpdCandles:stream.CloseAndRecv",
				zap.Int("sent nums", sentnums),
				zap.Error(err))
		} else {
			tradingdb2utils.Error("Client.UpdCandles:stream.CloseAndRecv",
				zap.Int("sent nums", sentnums),
				zap.Error(err))
		}

		// if error, reset
		client.reset()

		return nil, err
	}

	return reply, nil
}

// GetCandles - get candles
func (client *Client) GetCandles(ctx context.Context, market string, symbol string, tsStart int64, tsEnd int64, logger *zap.Logger) (
	*tradingpb.Candles, error) {

	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			if logger != nil {
				logger.Error("Client.GetCandles:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			} else {
				tradingdb2utils.Error("Client.GetCandles:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			}

			return nil, err
		}

		client.conn = conn
		client.client = tradingpb.NewTradingDB2Client(conn)
	}

	stream, err := client.client.GetCandles(ctx, &tradingpb.RequestGetCandles{
		Market:  market,
		Symbol:  symbol,
		TsStart: tsStart,
		TsEnd:   tsEnd,
		BasicRequest: &tradingpb.BasicRequestData{
			Token: client.token,
		},
	})
	if err != nil {
		if logger != nil {
			logger.Error("Client.GetCandles:client.GetCandles",
				zap.Error(err))
		} else {
			tradingdb2utils.Error("Client.GetCandles:client.GetCandles",
				zap.Error(err))
		}

		// if error, reset
		client.reset()

		return nil, err
	}

	candles := &tradingpb.Candles{
		Market: market,
		Symbol: symbol,
	}
	times := 0

	for {
		req, err := stream.Recv()
		if req != nil && (err == nil || err == io.EOF) {
			times++

			tradingdb2.MergeCandles(candles, req.Candles)
		}

		if err != nil {
			if err == io.EOF {
				return candles, nil
			}

			if logger != nil {
				logger.Error("Client.GetCandles:stream.Recv",
					zap.Int("times", times),
					zap.Int("candle nums", len(candles.Candles)),
					zap.Error(err))
			} else {
				tradingdb2utils.Error("Client.GetCandles:stream.Recv",
					zap.Int("times", times),
					zap.Int("candle nums", len(candles.Candles)),
					zap.Error(err))
			}

			// if error, reset
			client.reset()

			return nil, err
		}
	}
}

// UpdSymbol - update symbol
func (client *Client) UpdSymbol(ctx context.Context, si *tradingpb.SymbolInfo, logger *zap.Logger) (
	*tradingpb.ReplyUpdSymbol, error) {

	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			if logger != nil {
				logger.Error("Client.UpdSymbol:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			} else {
				tradingdb2utils.Error("Client.UpdSymbol:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			}

			return nil, err
		}

		client.conn = conn
		client.client = tradingpb.NewTradingDB2Client(conn)
	}

	reply, err := client.client.UpdSymbol(ctx, &tradingpb.RequestUpdSymbol{
		Symbol: si,
		BasicRequest: &tradingpb.BasicRequestData{
			Token: client.token,
		},
	})
	if err != nil {
		if logger != nil {
			logger.Error("Client.UpdSymbol:Client.UpdSymbol",
				zap.Error(err))
		} else {
			tradingdb2utils.Error("Client.UpdSymbol:Client.UpdSymbol",
				zap.Error(err))
		}

		// if error, reset
		client.reset()

		return nil, err
	}

	return reply, nil
}

// GetSymbol - get symbol
func (client *Client) GetSymbol(ctx context.Context, market string, symbol string, logger *zap.Logger) (
	*tradingpb.SymbolInfo, error) {

	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			if logger != nil {
				logger.Error("Client.GetSymbol:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			} else {
				tradingdb2utils.Error("Client.GetSymbol:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			}

			return nil, err
		}

		client.conn = conn
		client.client = tradingpb.NewTradingDB2Client(conn)
	}

	reply, err := client.client.GetSymbol(ctx, &tradingpb.RequestGetSymbol{
		Market: market,
		Symbol: symbol,
		BasicRequest: &tradingpb.BasicRequestData{
			Token: client.token,
		},
	})
	if err != nil {
		if logger != nil {
			logger.Error("Client.GetSymbol:client.GetSymbol",
				zap.Error(err))
		} else {
			tradingdb2utils.Error("Client.GetSymbol:client.GetSymbol",
				zap.Error(err))
		}

		// if error, reset
		client.reset()

		return nil, err
	}

	return reply.Symbol, nil
}

// GetSymbols - get symbols
func (client *Client) GetSymbols(ctx context.Context, market string, symbols []string, onGetSymbols FuncOnGetSymbols, logger *zap.Logger) error {

	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			if logger != nil {
				logger.Error("Client.GetSymbols:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			} else {
				tradingdb2utils.Error("Client.GetSymbols:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))
			}

			return err
		}

		client.conn = conn
		client.client = tradingpb.NewTradingDB2Client(conn)
	}

	stream, err := client.client.GetSymbols(ctx, &tradingpb.RequestGetSymbols{
		Market:  market,
		Symbols: symbols,
		BasicRequest: &tradingpb.BasicRequestData{
			Token: client.token,
		},
	})
	if err != nil {
		if logger != nil {
			logger.Error("Client.GetSymbols:client.GetSymbols",
				zap.Error(err))
		} else {
			tradingdb2utils.Error("Client.GetSymbols:client.GetSymbols",
				zap.Error(err))
		}

		// if error, reset
		client.reset()

		return err
	}

	for {
		req, err := stream.Recv()
		if req != nil && req.Symbol != nil {
			onGetSymbols(req.Symbol)
		}

		if err != nil {
			if err == io.EOF {
				return nil
			}

			if logger != nil {
				logger.Error("Client.GetSymbols:stream.Recv",
					zap.Error(err))
			} else {
				tradingdb2utils.Error("Client.GetSymbols:stream.Recv",
					zap.Error(err))
			}

			// if error, reset
			client.reset()

			return err
		}
	}
}
