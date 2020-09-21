package trdb2chatbot

import (
	"context"

	chatbot "github.com/zhs007/chatbot"
	chatbotbase "github.com/zhs007/chatbot/base"
	basicchatbot "github.com/zhs007/chatbot/basicchatbot"
	chatbotusermgr "github.com/zhs007/chatbot/usermgr"
	tradingdb2 "github.com/zhs007/tradingdb2"
	tradingdb2grpc "github.com/zhs007/tradingdb2/grpc"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	tradingdb2ver "github.com/zhs007/tradingdb2/ver"
	"go.uber.org/zap"
)

func startServ(ctx context.Context, cfg *tradingdb2.Config, endchan chan int) (
	*tradingdb2grpc.Serv, error) {

	serv, err := tradingdb2grpc.NewServ(cfg)
	if err != nil {
		tradingdb2utils.Info("start",
			zap.Error(err))

		return nil, err
	}

	tradingdb2utils.Info("startServ...",
		zap.String("version", tradingdb2ver.Version),
		zap.String("chatbot.version", chatbotbase.VERSION))

	go func() {
		err := serv.Start(ctx)
		if err != nil {
			tradingdb2utils.Warn("startServ:Start",
				zap.Error(err))
		}

		endchan <- 0
	}()

	return serv, nil
}

func startChatBot(ctx context.Context, chatbotcfg *chatbot.Config, trdb2serv *tradingdb2grpc.Serv, endchan chan int) error {

	chatbotbase.SetLogger(tradingdb2utils.GetLogger())

	mgr, err := chatbotusermgr.NewUserMgr(chatbotcfg.DBPath,
		"", chatbotcfg.DBEngine, nil)
	if err != nil {
		tradingdb2utils.Warn("startChatBot:NewUserMgr",
			zap.Error(err))

		return err
	}

	serv, err := chatbot.NewChatBotServ(chatbotcfg, mgr, &ServiceCore{})
	if err != nil {
		tradingdb2utils.Warn("startChatBot:NewChatBotServ",
			zap.Error(err))

		return err
	}

	// register file processor
	// serv.MgrFile.RegisterFileProcessor(&markdownFP{
	// 	servAda: servAda,
	// })
	// serv.MgrFile.RegisterFileProcessor(&excelFP{
	// 	servAda: servAda,
	// })

	serv.Init(context.Background())

	go func() {
		serv.Start(context.Background())
		if err != nil {
			tradingdb2utils.Warn("startChatBot:Start",
				zap.Error(err))
		}

		endchan <- 0
	}()

	return nil
}

// Start - start serv & chatbot
func Start(ctx context.Context, cfgfn string, chatbotcfgfn string) error {

	cfg, err := tradingdb2.LoadConfig(cfgfn)
	if err != nil {
		tradingdb2utils.Error("Start:LoadConfig",
			zap.Error(err))

		return err
	}

	tradingdb2utils.InitLogger("tradingdb2", tradingdb2ver.Version, cfg.LogLevel, true, cfg.LogPath)

	chatbotcfg, err := chatbot.LoadConfig(chatbotcfgfn)
	if err != nil {
		tradingdb2utils.Warn("Start:LoadConfig",
			zap.Error(err))

		return err
	}

	err = basicchatbot.InitBasicChatBot(*chatbotcfg)
	if err != nil {
		tradingdb2utils.Warn("Start:InitBasicChatBot",
			zap.Error(err))

		return err
	}

	servchan := make(chan int)

	serv, err := startServ(context.Background(), cfg, servchan)
	if err != nil {
		tradingdb2utils.Warn("Start:startServ",
			zap.Error(err))

		return err
	}

	err = startChatBot(context.Background(), chatbotcfg, serv, servchan)
	if err != nil {
		tradingdb2utils.Warn("StartChatBot:startChatBot",
			zap.Error(err))

		return err
	}

	endi := 0
	for {
		<-servchan
		endi++
		if endi >= 2 {
			break
		}
	}

	return nil
}
