package trdb2chatbot

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/zhs007/chatbot"
	chatbotbase "github.com/zhs007/chatbot/base"
	chatbotpb "github.com/zhs007/chatbot/chatbotpb"
	tradingdb2grpc "github.com/zhs007/tradingdb2/grpc"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
)

// nameCmd - name for command
const nameHistoryTasks = "historytasks"

// cmdHistoryTasks - historytasks help
type cmdHistoryTasks struct {
	serv *tradingdb2grpc.Serv
}

// RunCommand - run command
func (cmd *cmdHistoryTasks) RunCommand(ctx context.Context, serv *chatbot.Serv, params interface{},
	chat *chatbotpb.ChatMsg, ui *chatbotpb.UserInfo, ud proto.Message,
	scs chatbotpb.ChatBotService_SendChatServer) (bool, []*chatbotpb.ChatMsg, error) {

	if serv == nil {
		return true, nil, chatbotbase.ErrCmdInvalidServ
	}

	if serv.MgrText == nil {
		return true, nil, chatbotbase.ErrCmdInvalidServMgrText
	}

	var lst []*chatbotpb.ChatMsg
	str := "There are no task."

	arr := cmd.serv.TasksMgr.RecvHistory()
	if len(arr) > 0 {
		cs, err := chatbotbase.JSONFormat(arr)
		if err != nil {
			tradingdb2utils.Error("cmdHistoryTasks.JSONFormat",
				tradingdb2utils.JSON("arr", arr),
				zap.Error(err))

			return true, nil, err
		}

		str = cs
	}

	msgtxt := &chatbotpb.ChatMsg{
		Msg: str,
		Uai: chat.Uai,
	}

	lst = append(lst, msgtxt)

	return true, lst, nil
}

// ParseCommandLine - parse command line
func (cmd *cmdHistoryTasks) ParseCommandLine(cmdline []string, chat *chatbotpb.ChatMsg) (interface{}, error) {
	tradingdb2utils.Info("cmdHistoryTasks.ParseCommandLine")

	return nil, nil
}

// OnMessage - get message
func (cmd *cmdHistoryTasks) OnMessage(ctx context.Context, serv *chatbot.Serv, chat *chatbotpb.ChatMsg,
	ui *chatbotpb.UserInfo, ud proto.Message,
	scs chatbotpb.ChatBotService_SendChatServer) (bool, []*chatbotpb.ChatMsg, error) {

	return true, nil, chatbotbase.ErrCmdItsNotMine
}

// RegisterCmdHistoryTasks - register historytasks in command
func RegisterCmdHistoryTasks(serv *tradingdb2grpc.Serv) {
	chatbot.RegisterCommand(nameHistoryTasks, &cmdHistoryTasks{
		serv: serv,
	})
}
