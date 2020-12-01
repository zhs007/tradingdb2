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
	lis          net.Listener
	grpcServ     *grpc.Server
	DB           *tradingdb2.DB
	Cfg          *tradingdb2.Config
	DBSimTrading *tradingdb2.SimTradingDB
	MgrNodes     *Node2Mgr
}

// NewServ -
func NewServ(cfg *tradingdb2.Config) (*Serv, error) {

	db, err := tradingdb2.NewDB(cfg.DBPath, "", cfg.DBEngine)
	if err != nil {
		tradingdb2utils.Error("NewServ.NewDB",
			zap.Error(err))

		return nil, err
	}

	dbSimTrading, err := tradingdb2.NewSimTradingDB(cfg.DBPath, "", cfg.DBEngine)
	if err != nil {
		tradingdb2utils.Error("NewServ.NewSimTradingDB",
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

	mgrNodes, err := NewNode2Mgr(cfg)
	if err != nil {
		tradingdb2utils.Error("NewServ.NewNode2Mgr",
			zap.Error(err))

		return nil, err
	}

	serv := &Serv{
		lis:          lis,
		grpcServ:     grpcServ,
		DB:           db,
		DBSimTrading: dbSimTrading,
		Cfg:          cfg,
		MgrNodes:     mgrNodes,
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

// checkBasicRequest - check BasicRequest
func (serv *Serv) checkBasicRequest(req *tradingpb.BasicRequestData) error {
	if req.Token == "" || tradingdb2utils.IndexOfStringSlice(serv.Cfg.Tokens, req.Token, 0) < 0 {
		return tradingdb2.ErrInvalidToken
	}

	return nil
}

// UpdCandles - update candles
func (serv *Serv) UpdCandles(stream tradingpb.TradingDB2_UpdCandlesServer) error {
	candles := &tradingpb.Candles{}
	times := 0

	for {
		req, err := stream.Recv()
		if req != nil && (err == nil || err == io.EOF) {
			if times == 0 {
				err = serv.checkBasicRequest(req.BasicRequest)
				if err != nil {
					tradingdb2utils.Error("Serv.UpdCandles:Token",
						zap.Int("length", len(candles.Candles)),
						zap.Int("times", times),
						zap.String("token", req.Token),
						zap.Error(err))

					return err
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
	err := serv.checkBasicRequest(req.BasicRequest)
	if err != nil {
		tradingdb2utils.Error("Serv.GetCandles:checkToken",
			zap.String("token", req.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(err))

		return err
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
	err := serv.checkBasicRequest(req.BasicRequest)
	if err != nil {
		tradingdb2utils.Error("Serv.UpdSymbol:checkToken",
			zap.String("token", req.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(err))

		return nil, err
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
	err := serv.checkBasicRequest(req.BasicRequest)
	if err != nil {
		tradingdb2utils.Error("Serv.GetSymbol:checkToken",
			zap.String("token", req.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(err))

		return nil, err
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
	err := serv.checkBasicRequest(req.BasicRequest)
	if err != nil {
		tradingdb2utils.Error("Serv.GetSymbols:checkToken",
			zap.String("token", req.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(err))

		return err
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

// SimTrading - simTrading
func (serv *Serv) SimTrading(ctx context.Context, req *tradingpb.RequestSimTrading) (*tradingpb.ReplySimTrading, error) {
	err := serv.checkBasicRequest(req.BasicRequest)
	if err != nil {
		tradingdb2utils.Error("Serv.SimTrading:checkToken",
			zap.String("token", req.BasicRequest.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(err))

		return nil, err
	}

	params, err := serv.DB.FixSimTradingParams(ctx, req.Params)
	if err != nil {
		tradingdb2utils.Error("Serv.SimTrading:FixSimTradingParams",
			zap.Error(err))

		return nil, err
	}

	pnl, err := serv.DBSimTrading.GetSimTrading(ctx, params)
	if err != nil {
		tradingdb2utils.Error("Serv.SimTrading:GetSimTrading",
			zap.Error(err))

		return nil, err
	}

	if pnl != nil {
		return &tradingpb.ReplySimTrading{
			Pnl: []*tradingpb.PNLData{
				pnl,
			},
		}, nil
	}

	res := &tradingpb.ReplySimTrading{}

	for _, v := range params.Baselines {
		bp := GenCalcBaseline(v, params.StartTs, params.EndTs)
		br, err := serv.MgrNodes.CalcPNL(ctx, bp, nil)
		if err != nil {
			tradingdb2utils.Error("Serv.SimTrading:CalcPNL baseline",
				tradingdb2utils.JSON("parameter", bp),
				zap.Error(err))

			return nil, err
		}

		res.Baseline = append(res.Baseline, br.Pnl...)
	}

	reply, err := serv.MgrNodes.CalcPNL(ctx, params, nil)
	if err != nil {
		tradingdb2utils.Error("Serv.SimTrading:CalcPNL",
			zap.Error(err))

		return nil, err
	}

	// if len(reply.Pnl) > 0 {
	// 	for _, v := range reply.Pnl[0].Total.Values {
	// 		tradingdb2utils.Info("Serv.SimTrading:CalcPNL",
	// 			zap.Float32("perValue", v.PerValue),
	// 			zap.Float32("cost", v.Cost),
	// 			zap.Float32("value", v.Value),
	// 			zap.Int64("ts", v.Ts))
	// 	}
	// }

	res.Pnl = reply.Pnl

	if len(reply.Pnl) > 0 {
		err = serv.DBSimTrading.UpdSimTrading(ctx, params, reply.Pnl[0])
		if err != nil {
			tradingdb2utils.Error("Serv.SimTrading:UpdSimTrading",
				zap.Error(err))

			return nil, err
		}
	}

	return res, nil
}
