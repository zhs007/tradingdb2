package main

import (
	"context"
	"fmt"

	trdb2chatbot "github.com/zhs007/tradingdb2/chatbot"
)

func main() {
	//! 这里之所以用 config 目录，是因为docker映射目录时，可以映射config目录，
	//! 这样cfg里面的配置，可以完全保留
	err := trdb2chatbot.Start(context.Background(), "./cfg/config.yaml", "./cfg/chatbot.yaml")
	if err != nil {
		fmt.Printf("StartChatBot %v", err)

		return
	}
}
