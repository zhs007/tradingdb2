package tradingdb2grpc

import (
	"context"
	"io"
	"net"

	tradingdb2 "github.com/zhs007/tradingdb2"
	tradingdb2pb "github.com/zhs007/tradingdb2/tradingdb2pb"
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

	tradingdb2pb.RegisterTradingDB2ServiceServer(grpcServ, serv)

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
func (serv *Serv) UpdCandles(stream tradingdb2pb.TradingDB2Service_UpdCandlesServer) error {
	candles := &tradingdb2pb.Candles{}
	times := 0
	// token := ""

	for {
		req, err := stream.Recv()
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

				return stream.SendAndClose(&tradingdb2pb.ReplyUpdCandles{
					LengthOK: int32(len(candles.Candles)),
				})
			}

			tradingdb2utils.Error("Serv.UpdCandles:Recv",
				zap.Int("length", len(candles.Candles)),
				zap.Int("times", times),
				zap.Error(err))

			return err
		}

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

			// token = req.Token

			candles.Market = req.Candles.Market
			candles.Symbol = req.Candles.Symbol
			candles.Tag = req.Candles.Tag
		}
		// } else {
		// 	if token != req.Token {
		// 		tradingdb2utils.Error("Serv.UpdCandles:Token",
		// 			zap.Int("length", len(candles.Candles)),
		// 			zap.Int("times", times),
		// 			zap.String("first token", token),
		// 			zap.String("token", req.Token),
		// 			zap.Error(tradingdb2.ErrInvalidToken))

		// 		return tradingdb2.ErrInvalidToken
		// 	}

		// 	if candles.Market != req.Candles.Market {
		// 		tradingdb2utils.Error("Serv.UpdCandles:Market",
		// 			zap.Int("length", len(candles.Candles)),
		// 			zap.Int("times", times),
		// 			zap.String("first market", candles.Market),
		// 			zap.String("market", req.Candles.Market),
		// 			zap.Error(tradingdb2.ErrInvalidMarket))

		// 		return tradingdb2.ErrInvalidMarket
		// 	}

		// 	if candles.Symbol != req.Candles.Symbol {
		// 		tradingdb2utils.Error("Serv.UpdCandles:Symbol",
		// 			zap.Int("length", len(candles.Candles)),
		// 			zap.Int("times", times),
		// 			zap.String("first symbol", candles.Symbol),
		// 			zap.String("symbol", req.Candles.Symbol),
		// 			zap.Error(tradingdb2.ErrInvalidSymbol))

		// 		return tradingdb2.ErrInvalidSymbol
		// 	}

		// 	if candles.Tag != req.Candles.Tag {
		// 		tradingdb2utils.Error("Serv.UpdCandles:Tag",
		// 			zap.Int("length", len(candles.Candles)),
		// 			zap.Int("times", times),
		// 			zap.String("first tag", candles.Tag),
		// 			zap.String("tag", req.Candles.Tag),
		// 			zap.Error(tradingdb2.ErrInvalidTag))

		// 		return tradingdb2.ErrInvalidTag
		// 	}
		// }

		times++

		tradingdb2.MergeCandles(candles, req.Candles)
	}
}

// GetCandles - get candles
func (serv *Serv) GetCandles(req *tradingdb2pb.RequestGetCandles, stream tradingdb2pb.TradingDB2Service_GetCandlesServer) error {
	if req.Token == "" || tradingdb2utils.IndexOfStringSlice(serv.Cfg.Tokens, req.Token, 0) < 0 {
		tradingdb2utils.Error("Serv.UpdCandles:checkToken",
			zap.String("token", req.Token),
			zap.Strings("tokens", serv.Cfg.Tokens),
			zap.Error(tradingdb2.ErrInvalidToken))

		return tradingdb2.ErrInvalidToken
	}

	candles, err := serv.DB.GetCandles(stream.Context(), req.Market, req.Symbol, req.Tag)
	if err != nil {
		tradingdb2utils.Error("Serv.UpdCandles:DB.UpdCandles",
			tradingdb2utils.JSON("params", req),
			zap.Error(err))

		return err
	}

	// sentnums := 0

	return tradingdb2.BatchCandles(candles, serv.Cfg.BatchCandleNums, func(lst *tradingdb2pb.Candles) error {
		cc := &tradingdb2pb.ReplyGetCandles{
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
