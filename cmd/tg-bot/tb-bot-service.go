package main

import (
	"fmt"
	"github.com/digital-technology-agency/secret-santa/pkg/models"
	"github.com/digital-technology-agency/secret-santa/pkg/services"
	"github.com/digital-technology-agency/secret-santa/pkg/utils"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

var (
	address  = utils.GetEnv("PORT", "8080")
	rooms    = map[string]*services.Game{}
	keyboard = tgbot.NewInlineKeyboardMarkup(
		tgbot.NewInlineKeyboardRow(
			tgbot.NewInlineKeyboardButtonData("🎅🏻 Список участников", services.CmdLayerListGame),
			tgbot.NewInlineKeyboardButtonData("🎄 Хочу в игру", services.CmdJoinGame),
			tgbot.NewInlineKeyboardButtonData("❌ Выход из игры", services.CmdExitGame),
			tgbot.NewInlineKeyboardButtonData("🌐 Выбор языка", services.CmdLanguageGame),
		),
	)
)

func initRoom(id string) *services.Game {
	if rooms[id] == nil {
		create, err := services.GetOrCreate(id)
		if err != nil {
			log.Panic(err)
		}
		rooms[id] = create
	}
	return rooms[id]
}

// main
func main() {
	/*	rt := router.New()
		rt.GET("/", routes.GetHealth)
		webServer := fasthttp.Server{
			Name:         "Santa bot",
			WriteTimeout: time.Second * 5,
			ReadTimeout:  time.Second * 5,
			IdleTimeout:  time.Second * 5,
			Handler:      rt.Handler,
		}
		go func() {
			fmt.Print("GET... [http://localhost", address, "/", "]\n")
			if err := webServer.ListenAndServe(fmt.Sprintf(":%s", address)); err != nil {
				log.Panic(err)
			}
		}()*/
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
			chatId := fmt.Sprintf("%d", update.CallbackQuery.Message.Chat.ID)
			userId := fmt.Sprintf("%d", update.CallbackQuery.From.ID)
			lastName := update.CallbackQuery.From.FirstName
			msgConfig = tgbot.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Пользователь %s - функция находиться в разработке!", lastName))
			game := initRoom(chatId)
			switch cmd {
			default:
				log.Println("Cmd:[%s] - not found", cmd)
			case services.CmdLayerListGame:
				msgConfig = tgbot.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Пользователь %s - запросил список участников игры!", lastName))
				players, _ := game.GetAllPlayers()
				btns := tgbot.NewInlineKeyboardRow()
				for _, player := range players {
					btns = append(btns, tgbot.NewInlineKeyboardButtonData(fmt.Sprintf("🎅🏻 %s", player.Login), player.Login))
				}
				if len(players) > 0 {
					msgConfig.ReplyMarkup = tgbot.NewInlineKeyboardMarkup(btns)
				}
			case services.CmdJoinGame:
				player, _ := game.GetPlayer(userId)
				if player == nil {
					err = game.AddPlayer(models.Player{
						Id:       userId,
						Login:    lastName,
						FriendId: "",
					})
					if err != nil {
						log.Panic(err)
					}
					err = game.Algorithm()
					if err != nil {
						log.Panic(err)
					}
					msgConfig = tgbot.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Пользователь %s - присоединился к игре!", lastName))
				} else {
					msgConfig = tgbot.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Пользователь %s - уже является участником игры!", lastName))
				}
			case services.CmdExitGame:
				err := game.RemovePlayerById(userId)
				if err != nil {
					log.Panic(err)
				}
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
