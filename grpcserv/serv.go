package tradingdb2serv

import (
	"context"
	"net"

	tradingdb2 "github.com/zhs007/tradingdb2"
	tradingdb2pb "github.com/zhs007/tradingdb2/pb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	tradingdb2ver "github.com/zhs007/tradingdb2/ver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Serv - tradingdb2 Service
type Serv struct {
	lis      net.Listener
	grpcServ *grpc.Server
	db       *tradingdb2.DB
	cfg      *tradingdb2.Config
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
		db:       db,
		cfg:      cfg,
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
func (serv *Serv) UpdCandles(ctx context.Context, req *tradingdb2pb.RequestUpdCandles) (*tradingdb2pb.ReplyUpdCandles, error) {
	return nil, nil
}
