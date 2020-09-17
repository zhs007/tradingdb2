package main

import (
	"context"

	tradingdb2 "github.com/zhs007/tradingdb2"
	tradingdb2serv "github.com/zhs007/tradingdb2/grpcserv"
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

	serv, err := tradingdb2serv.NewServ(cfg)
	if err != nil {
		tradingdb2utils.Error("NewServ",
			zap.Error(err))

		return
	}

	serv.Start(context.Background())
}
