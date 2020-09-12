package tradingdb2serv

import (
	"context"
	"net"

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
}

// NewServ -
func NewServ(cfgfn string) (*Serv, error) {
	cfg, err := LoadConfig(cfgfn)
	if err != nil {
		tradingdb2utils.Error("tradingdb2serv.NewServ",
			zap.String("cfgfn", cfgfn),
			zap.Error(err))
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

// UpdCandle - update candle
func (serv *Serv) UpdCandle(context.Context, *tradingdb2pb.Candle) (*tradingdb2pb.ReplyUpdCandle, error) {
	return nil, nil
}
