package main

import (
	"fmt"
	"github.com/digital-technology-agency/secret-santa/pkg/services"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

var (
	keyboard = tgbot.NewInlineKeyboardMarkup(
		tgbot.NewInlineKeyboardRow(
			tgbot.NewInlineKeyboardButtonData("🎅🏻 Список участников", services.CmdLayerListGame),
			tgbot.NewInlineKeyboardButtonData("🎄 Хочу в игру", services.CmdJoinGame),
			tgbot.NewInlineKeyboardButtonData("❌ Выход из игры", services.CmdExitGame),
			tgbot.NewInlineKeyboardButtonData("🌐 Выбор языка", services.CmdLanguageGame),
		),
	)
)

func main() {
	bot, err := tgbot.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	update := tgbot.NewUpdate(0)
	update.Timeout = 60
	updatesChan := bot.GetUpdatesChan(update)
	for update := range updatesChan {
		if update.Message != nil {
			msgText := update.Message.Text
			if services.InitGameRegex.MatchString(msgText) {
				msg := tgbot.NewMessage(update.Message.Chat.ID, "Вы может принять участие в игре")
				msg.ReplyMarkup = keyboard
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}
		} else if update.CallbackQuery != nil {
			var msgConfig tgbot.MessageConfig
			cmd := update.CallbackQuery.Data
			lastName := update.CallbackQuery.Message.Chat.FirstName
			switch cmd {
			default:
				log.Panicf("Cmd:[%s] - not found", cmd)
			case services.CmdLayerListGame:
				msgConfig = tgbot.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Пользователь %s - запросил список участников игры!", lastName))
			case services.CmdJoinGame:
				msgConfig = tgbot.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Пользователь %s - присоединился к игре!", lastName))
			case services.CmdExitGame:
				msgConfig = tgbot.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Пользователь %s - вышел из игры!", lastName))
			case services.CmdLanguageGame:
				msgConfig = tgbot.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Пользователь %s - изменил язык игры!", lastName))
			}
			if msgConfig.Text == "" {
				continue
			}
			if _, err := bot.Send(msgConfig); err != nil {
				log.Panic(err)
			}
		}
	}
}
