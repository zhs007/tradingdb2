package tradingdb2task

import (
	"encoding/hex"

	"github.com/zhs007/tradingdb2/tradingpb"
	tradingdb2utils "github.com/zhs007/tradingdb2/utils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type TaskKeyList struct {
	lst []string
}

func NewTaskKeyList() *TaskKeyList {
	return &TaskKeyList{}
}

func (lst *TaskKeyList) AddTask(params *tradingpb.SimTradingParams) error {
	buf, err := proto.Marshal(params)
	if err != nil {
		tradingdb2utils.Warn("TaskKeyList.AddTask:Marshal",
			zap.Error(err))

		return err
	}

	lst.lst = append(lst.lst, hex.EncodeToString(buf))

	return nil
}

func (lst *TaskKeyList) RemoveTask(task []byte) error {
	key := hex.EncodeToString(task)

	for i, v := range lst.lst {
		if key == v {
			lst.lst = append(lst.lst[:i], lst.lst[i+1:]...)

			return nil
		}
	}

	tradingdb2utils.Warn("TaskKeyList:RemoveTask",
		zap.String("key", key),
		zap.Error(ErrNoKey))

	return ErrNoKey
}
