package tradingdb2grpc

import (
	"context"
	"io"
	"net"

	tradingdb2 "github.com/zhs007/tradingdb2"
	tradingpb "github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	tradingdb2ver "github.com/zhs007/tradingdb2/ver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Serv - tradingdb2 Service
type Serv struct {
	lis      net.Listener
	grpcServ *grpc.Server
	DB       *tradingdb2.DB
	Cfg      *tradingdb2.Config
}

// NewServ -
func NewServ(cfg *tradingdb2.Config) (*Serv, error) {

	db, err := tradingdb2.NewDB(cfg.DBPath, "", cfg.DBEngine)
	if err != nil {
		tradingdb2utils.Error("NewServ.NewDB",
			zap.Error(err))

		return nil, err
	}

	lis, err := net.Listen("tcp", cfg.BindAddr)
	if err != nil {
		tradingdb2utils.Error("NewServ.Listen",
			zap.Error(err))

		return nil, err
	}

	grpcServ := grpc.NewServer()

	serv := &Serv{
		lis:      lis,
		grpcServ: grpcServ,
		DB:       db,
		Cfg:      cfg,
	}

	tradingpb.RegisterTradingDB2Server(grpcServ, serv)

	tradingdb2utils.Info("NewServ OK.",
		zap.String("addr", cfg.BindAddr),
		zap.String("ver", tradingdb2ver.Version))

	return serv, nil
}

// Start - start a service
func (serv *Serv) Start(ctx context.Context) error {
	return serv.grpcServ.Serve(serv.lis)
}

// Stop - stop service
func (serv *Serv) Stop() {
	serv.lis.Close()

	return
}

// UpdCandles - update candles
func (serv *Serv) UpdCandles(stream tradingpb.TradingDB2_UpdCandlesServer) error {
	candles := &tradingpb.Candles{}
	times := 0
	// token := ""

	for {
		req, err := stream.Recv()
		if req != nil && (err == nil || err == io.EOF) {
			if times == 0 {
				if req.Token == "" || tradingdb2utils.IndexOfStringSlice(serv.Cfg.Tokens, req.Token, 0) < 0 {
					tradingdb2utils.Error("Serv.UpdCandles:Token",
						zap.Int("length", len(candles.Candles)),
						zap.Int("times", times),
						zap.String("token", req.Token),
						zap.Error(tradingdb2.ErrInvalidToken))

					return tradingdb2.ErrInvalidToken
				}

				if req.Candles.Market == "" ||
					req.Candles.Symbol == "" ||
					req.Candles.Tag == "" {

					tradingdb2utils.Error("Serv.UpdCandles:Market|Symbol|Tag",
						zap.Int("length", len(candles.Candles)),
						zap.Int("times", times),
						zap.String("market", req.Candles.Market),
						zap.String("symbol", req.Candles.Symbol),
						zap.String("tag", req.Candles.Tag),
						zap.Error(tradingdb2.ErrInvalidUpdCandlesParams))

					return tradingdb2.ErrInvalidUpdCandlesParams
				}

				candles.Market = req.Candles.Market
				candles.Symbol = req.Candles.Symbol
				candles.Tag = req.Candles.Tag
			}

			times++

			tradingdb2.MergeCandles(candles, req.Candles)
		}

		if err != nil {
			if err == io.EOF {
				err := serv.DB.UpdCandles(stream.Context(), candles)
				if err != nil {
					tradingdb2utils.Error("Serv.UpdCandles:DB.UpdCandles",
						zap.Int("length", len(candles.Candles)),
						zap.Int("times", times),
						zap.Error(err))

					return err
				}

				return stream.SendAndClose(&tradingpb.ReplyUpdCandles{
					LengthOK: int32(len(candles.Candles)),
				})
			}

			tradingdb2utils.Error("Serv.UpdCandles:Recv",
				zap.Int("length", len(candles.Candles)),
				zap.Int("times", times),
				zap.Error(err))

			return err
		}
	}
}

// GetCandles - get candles
func (serv *Serv) GetCandles(req *tradingpb.RequestGetCandles, stream tradingpb.TradingDB2_GetCandlesServer) error {
	if req.Token == "" || tradingdb2utils.IndexOfStringSlice(serv.Cfg.Tokens, req.Token, 0) < 0 {
		tradingdb2utils.Error("Serv.GetCandles:checkToken",
			zap.String("token", req.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(tradingdb2.ErrInvalidToken))

		return tradingdb2.ErrInvalidToken
	}

	var tags []string
	if len(req.Tags) > 0 {
		tags = req.Tags
	} else if len(req.Tags) == 0 && req.Tag != "" {
		tags = append(tags, req.Tag)
	}

	candles, err := serv.DB.GetCandles(stream.Context(), req.Market, req.Symbol, tags, req.TsStart, req.TsEnd)
	if err != nil {
		tradingdb2utils.Error("Serv.GetCandles:DB.GetCandles",
			tradingdb2utils.JSON("params", req),
			zap.Error(err))

		return err
	}

	// sentnums := 0

	return tradingdb2.BatchCandles(candles, serv.Cfg.BatchCandleNums, func(lst *tradingpb.Candles) error {
		cc := &tradingpb.ReplyGetCandles{
			Candles: lst,
		}

		// if sentnums == 0 {
		// 	cc.Candles.Market = candles.Market
		// 	cc.Candles.Symbol = candles.Symbol
		// 	cc.Candles.Tag = candles.Tag
		// }

		return stream.Send(cc)
	})
}

// UpdSymbol - update symbol
func (serv *Serv) UpdSymbol(ctx context.Context, req *tradingpb.RequestUpdSymbol) (*tradingpb.ReplyUpdSymbol, error) {
	if req.Token == "" || tradingdb2utils.IndexOfStringSlice(serv.Cfg.Tokens, req.Token, 0) < 0 {
		tradingdb2utils.Error("Serv.UpdSymbol:checkToken",
			zap.String("token", req.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(tradingdb2.ErrInvalidToken))

		return nil, tradingdb2.ErrInvalidToken
	}

	symbol, err := serv.DB.GetSymbol(ctx, req.Symbol.Market, req.Symbol.Symbol)
	if err != nil {
		tradingdb2utils.Error("Serv.UpdSymbol:DB.GetSymbol",
			zap.Error(err))

		return nil, err
	}

	if symbol != nil {
		symbol.Fund = tradingdb2.MergeFund(symbol.Fund, req.Symbol.Fund)
	} else {
		symbol = req.Symbol
	}

	err = serv.DB.UpdSymbol(ctx, symbol)
	if err != nil {
		tradingdb2utils.Error("Serv.UpdSymbol:DB.UpdSymbol",
			zap.Error(err))

		return nil, err
	}

	res := &tradingpb.ReplyUpdSymbol{
		IsOK: true,
	}

	return res, nil
}

// GetSymbol - get symbol
func (serv *Serv) GetSymbol(ctx context.Context, req *tradingpb.RequestGetSymbol) (*tradingpb.ReplyGetSymbol, error) {
	if req.Token == "" || tradingdb2utils.IndexOfStringSlice(serv.Cfg.Tokens, req.Token, 0) < 0 {
		tradingdb2utils.Error("Serv.GetSymbol:checkToken",
			zap.String("token", req.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(tradingdb2.ErrInvalidToken))

		return nil, tradingdb2.ErrInvalidToken
	}

	si, err := serv.DB.GetSymbol(ctx, req.Market, req.Symbol)
	if err != nil {
		tradingdb2utils.Error("Serv.GetSymbol:DB.GetSymbol",
			zap.String("market", req.Market),
			zap.String("sytmbol", req.Symbol),
			zap.Error(err))

		return nil, err
	}

	res := &tradingpb.ReplyGetSymbol{
		Symbol: si,
	}

	return res, nil
}

// GetSymbols - get symbols
func (serv *Serv) GetSymbols(req *tradingpb.RequestGetSymbols, stream tradingpb.TradingDB2_GetSymbolsServer) error {
	if req.Token == "" || tradingdb2utils.IndexOfStringSlice(serv.Cfg.Tokens, req.Token, 0) < 0 {
		tradingdb2utils.Error("Serv.GetSymbols:checkToken",
			zap.String("token", req.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(tradingdb2.ErrInvalidToken))

		return tradingdb2.ErrInvalidToken
	}

	var symbols []string
	if len(req.Symbols) > 0 {
		symbols = req.Symbols
	} else {
		arr, err := serv.DB.GetMarketSymbols(stream.Context(), req.Market)
		if err != nil {
			tradingdb2utils.Error("Serv.GetSymbols:GetMarketSymbols",
				zap.String("Market", req.Market),
				zap.Error(err))

			return err
		}

		symbols = arr
	}

	for _, v := range symbols {
		si, err := serv.DB.GetSymbol(stream.Context(), req.Market, v)
		if err != nil {
			tradingdb2utils.Error("Serv.GetSymbols:DB.GetSymbol",
				zap.String("market", req.Market),
				zap.String("sytmbol", v),
				zap.Error(err))

			// return err
		}

		if si != nil {
			res := &tradingpb.ReplyGetSymbol{
				Symbol: si,
			}

			err = stream.Send(res)
			if err != nil {
				tradingdb2utils.Error("Serv.GetSymbols:Send",
					zap.String("sytmbol", v),
					zap.Error(err))

				// return err
			}
		}
	}

	return nil
}
