package trdb2chatbot

import (
	"context"
	"path"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/pflag"
	"github.com/zhs007/chatbot"
	chatbotbase "github.com/zhs007/chatbot/base"
	chatbotpb "github.com/zhs007/chatbot/chatbotpb"
	tradingdb2 "github.com/zhs007/tradingdb2"
	tradingdb2grpc "github.com/zhs007/tradingdb2/grpc"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
)

// nameCmd - name for command
const nameExpSymbolsCmd = "expsymbols"

type paramsCmdExpSymbols struct {
	market string
}

// cmdExpSymbols - expsymbols help
type cmdExpSymbols struct {
	serv *tradingdb2grpc.Serv
}

// RunCommand - run command
func (cmd *cmdExpSymbols) RunCommand(ctx context.Context, serv *chatbot.Serv, params interface{},
	chat *chatbotpb.ChatMsg, ui *chatbotpb.UserInfo, ud proto.Message,
	scs chatbotpb.ChatBotService_SendChatServer) (bool, []*chatbotpb.ChatMsg, error) {

	if serv == nil {
		return true, nil, chatbotbase.ErrCmdInvalidServ
	}

	if serv.MgrText == nil {
		return true, nil, chatbotbase.ErrCmdInvalidServMgrText
	}

	cmdes, isok := params.(paramsCmdExpSymbols)
	if !isok {
		return true, nil, ErrCmdInvalidParams
	}

	fn := tradingdb2utils.AppendString(cmdes.market, ".", time.Now().UTC().Format("2006-01-02_15:04:05"), ".xlsx")

	err := tradingdb2.ExpSymbols(ctx, path.Join(cmd.serv.Cfg.DataPath, fn), cmd.serv.DB, cmdes.market)
	if err != nil {
		tradingdb2utils.Error("cmdExpSymbols.ExpSymbols",
			zap.Error(err))

		return true, nil, err
	}

	var lst []*chatbotpb.ChatMsg
	msgtxt := &chatbotpb.ChatMsg{
		Msg: tradingdb2utils.AppendString(cmd.serv.Cfg.DataURL, fn),
		Uai: chat.Uai,
	}

	lst = append(lst, msgtxt)

	return true, lst, nil
}

// ParseCommandLine - parse command line
func (cmd *cmdExpSymbols) ParseCommandLine(cmdline []string, chat *chatbotpb.ChatMsg) (interface{}, error) {

	flagset := pflag.NewFlagSet(nameExpSymbolsCmd, pflag.ContinueOnError)

	strMarket := flagset.StringP("market", "m", "", "market")

	err := flagset.Parse(cmdline[1:])
	if err != nil {
		return nil, err
	}

	return paramsCmdExpSymbols{
		market: *strMarket,
	}, nil
}

// OnMessage - get message
func (cmd *cmdExpSymbols) OnMessage(ctx context.Context, serv *chatbot.Serv, chat *chatbotpb.ChatMsg,
	ui *chatbotpb.UserInfo, ud proto.Message,
	scs chatbotpb.ChatBotService_SendChatServer) (bool, []*chatbotpb.ChatMsg, error) {

	return true, nil, chatbotbase.ErrCmdItsNotMine
}

// RegistercmdExpSymbols - register expalldata in command
func RegistercmdExpSymbols(serv *tradingdb2grpc.Serv) {
	chatbot.RegisterCommand(nameExpAllDataCmd, &cmdExpSymbols{
		serv: serv,
	})
}
