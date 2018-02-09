package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var botRef *tgbotapi.BotAPI = nil

func botNotify(message string) bool {
	if botRef == nil {
		bt, _ := tgbotapi.NewBotAPI(configuration.BotToken)
		botRef = bt
	}
	msg := tgbotapi.NewMessageToChannel(configuration.ChannelName, message)
	_, err := botRef.Send(msg)
	return err == nil
}
