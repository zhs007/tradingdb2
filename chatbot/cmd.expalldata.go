package trdb2chatbot

import (
	"context"
	"path"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/zhs007/chatbot"
	chatbotbase "github.com/zhs007/chatbot/base"
	chatbotpb "github.com/zhs007/chatbot/chatbotpb"
	tradingdb2 "github.com/zhs007/tradingdb2"
	tradingdb2grpc "github.com/zhs007/tradingdb2/grpc"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
)

// nameCmd - name for command
const nameExpAllDataCmd = "expalldata"

// cmdExpAllData - expalldata help
type cmdExpAllData struct {
	serv *tradingdb2grpc.Serv
}

// RunCommand - run command
func (cmd *cmdExpAllData) RunCommand(ctx context.Context, serv *chatbot.Serv, params interface{},
	chat *chatbotpb.ChatMsg, ui *chatbotpb.UserInfo, ud proto.Message,
	scs chatbotpb.ChatBotService_SendChatServer) (bool, []*chatbotpb.ChatMsg, error) {

	if serv == nil {
		return true, nil, chatbotbase.ErrCmdInvalidServ
	}

	if serv.MgrText == nil {
		return true, nil, chatbotbase.ErrCmdInvalidServMgrText
	}

	root, err := cmd.serv.DB.GetAllData(ctx)
	if err != nil {
		tradingdb2utils.Error("cmdExpAllData.RunCommand",
			zap.Error(err))

		return true, nil, err
	}

	fn := tradingdb2utils.AppendString("alldata.", time.Now().UTC().Format("2006-01-02_15:04:05"), ".xlsx")
	tradingdb2.ExpAllData(path.Join(cmd.serv.Cfg.DataPath, fn), root)

	var lst []*chatbotpb.ChatMsg
	msgtxt := &chatbotpb.ChatMsg{
		Msg: tradingdb2utils.AppendString(cmd.serv.Cfg.DataURL, fn),
		Uai: chat.Uai,
	}

	lst = append(lst, msgtxt)

	return true, lst, nil
}

// ParseCommandLine - parse command line
func (cmd *cmdExpAllData) ParseCommandLine(cmdline []string, chat *chatbotpb.ChatMsg) (interface{}, error) {

	return nil, nil
}

// OnMessage - get message
func (cmd *cmdExpAllData) OnMessage(ctx context.Context, serv *chatbot.Serv, chat *chatbotpb.ChatMsg,
	ui *chatbotpb.UserInfo, ud proto.Message,
	scs chatbotpb.ChatBotService_SendChatServer) (bool, []*chatbotpb.ChatMsg, error) {

	return true, nil, chatbotbase.ErrCmdItsNotMine
}

// RegisterCmdExpAllData - register expalldata in command
func RegisterCmdExpAllData(serv *tradingdb2grpc.Serv) {
	chatbot.RegisterCommand(nameExpAllDataCmd, &cmdExpAllData{
		serv: serv,
	})
}
