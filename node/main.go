package main

import (
	"context"

	tradingdb2 "github.com/zhs007/tradingdb2"
	tradingdb2grpc "github.com/zhs007/tradingdb2/grpc"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	tradingdb2ver "github.com/zhs007/tradingdb2/ver"
	"go.uber.org/zap"
)

func main() {
	cfg, err := tradingdb2.LoadConfig("./cfg/config.yaml")
	if err != nil {
		tradingdb2utils.Error("LoadConfig",
			zap.Error(err))

		return
	}

	tradingdb2utils.InitLogger("tradingdb2", tradingdb2ver.Version, cfg.LogLevel, true, "./logs")

	serv, err := tradingdb2grpc.NewServ(cfg)
	if err != nil {
		tradingdb2utils.Error("NewServ",
			zap.Error(err))

		return
	}

	tradingdb2utils.Error("Start TradingDB2 ...",
		zap.String("version", tradingdb2ver.Version),
		zap.String("bind address", cfg.BindAddr))

	serv.Start(context.Background())
}
