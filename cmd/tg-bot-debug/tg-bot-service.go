package main

import (
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var numericKeyboard = tgbot.NewReplyKeyboard(
	tgbot.NewKeyboardButtonRow(
		tgbot.NewKeyboardButtonContact("1"),
	),
)

func main() {
	bot, err := tgbot.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbot.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}
		msg := tgbot.NewMessage(update.Message.Chat.ID, update.Message.Text)
		switch update.Message.Text {
		case "хочу играть":
			msg.ReplyMarkup = numericKeyboard
		case "close":
			msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
		}
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
