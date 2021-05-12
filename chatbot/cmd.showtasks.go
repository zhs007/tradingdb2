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
const nameShowTasks = "showtasks"

// cmdShowTasks - expalldata help
type cmdShowTasks struct {
	serv *tradingdb2grpc.Serv
}

// RunCommand - run command
func (cmd *cmdShowTasks) RunCommand(ctx context.Context, serv *chatbot.Serv, params interface{},
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

	arr := cmd.serv.TasksMgr.GetTaskGroups()
	if len(arr) > 0 {
		cs, err := chatbotbase.JSONFormat(arr)
		if err != nil {
			tradingdb2utils.Error("cmdShowTasks.JSONFormat",
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
func (cmd *cmdShowTasks) ParseCommandLine(cmdline []string, chat *chatbotpb.ChatMsg) (interface{}, error) {
	tradingdb2utils.Info("cmdShowTasks.ParseCommandLine")

	return nil, nil
}

// OnMessage - get message
func (cmd *cmdShowTasks) OnMessage(ctx context.Context, serv *chatbot.Serv, chat *chatbotpb.ChatMsg,
	ui *chatbotpb.UserInfo, ud proto.Message,
	scs chatbotpb.ChatBotService_SendChatServer) (bool, []*chatbotpb.ChatMsg, error) {

	return true, nil, chatbotbase.ErrCmdItsNotMine
}

// RegisterCmdShowTasks - register showtasks in command
func RegisterCmdShowTasks(serv *tradingdb2grpc.Serv) {
	chatbot.RegisterCommand(nameShowTasks, &cmdShowTasks{
		serv: serv,
	})
}
