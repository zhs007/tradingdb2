package trdb2chatbot

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/pflag"
	"github.com/zhs007/chatbot"
	chatbotbase "github.com/zhs007/chatbot/base"
	chatbotpb "github.com/zhs007/chatbot/chatbotpb"
	tradingdb2grpc "github.com/zhs007/tradingdb2/grpc"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
)

// nameCmd - name for command
const nameShowLastTasks = "showlasttasks"

type paramsCmdShowLastTasks struct {
	taskGroupID int
}

// cmdShowLastTasks - showlasttasks help
type cmdShowLastTasks struct {
	serv *tradingdb2grpc.Serv
}

// RunCommand - run command
func (cmd *cmdShowLastTasks) RunCommand(ctx context.Context, serv *chatbot.Serv, params interface{},
	chat *chatbotpb.ChatMsg, ui *chatbotpb.UserInfo, ud proto.Message,
	scs chatbotpb.ChatBotService_SendChatServer) (bool, []*chatbotpb.ChatMsg, error) {

	if serv == nil {
		return true, nil, chatbotbase.ErrCmdInvalidServ
	}

	if serv.MgrText == nil {
		return true, nil, chatbotbase.ErrCmdInvalidServMgrText
	}

	cmdslt, isok := params.(paramsCmdShowLastTasks)
	if !isok {
		return true, nil, ErrCmdInvalidParams
	}

	var lst []*chatbotpb.ChatMsg
	str := "There are no task."

	arr := cmd.serv.TasksMgr.GetLastTasks(cmdslt.taskGroupID)
	if len(arr) > 0 {
		cs, err := chatbotbase.JSONFormat(arr)
		if err != nil {
			tradingdb2utils.Error("cmdShowLastTasks.JSONFormat",
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
func (cmd *cmdShowLastTasks) ParseCommandLine(cmdline []string, chat *chatbotpb.ChatMsg) (interface{}, error) {
	flagset := pflag.NewFlagSet(nameShowLastTasks, pflag.ContinueOnError)

	taskGroupID := flagset.IntP("taskgroup", "g", -1, "taskgroup")

	err := flagset.Parse(cmdline[1:])
	if err != nil {
		return nil, err
	}

	tradingdb2utils.Info("cmdShowLastTasks.ParseCommandLine",
		zap.Int("taskgroup", *taskGroupID))

	return paramsCmdShowLastTasks{
		taskGroupID: *taskGroupID,
	}, nil
}

// OnMessage - get message
func (cmd *cmdShowLastTasks) OnMessage(ctx context.Context, serv *chatbot.Serv, chat *chatbotpb.ChatMsg,
	ui *chatbotpb.UserInfo, ud proto.Message,
	scs chatbotpb.ChatBotService_SendChatServer) (bool, []*chatbotpb.ChatMsg, error) {

	return true, nil, chatbotbase.ErrCmdItsNotMine
}

// RegisterCmdShowLastTasks - register showlasttasks in command
func RegisterCmdShowLastTasks(serv *tradingdb2grpc.Serv) {
	chatbot.RegisterCommand(nameShowLastTasks, &cmdShowLastTasks{
		serv: serv,
	})
}
