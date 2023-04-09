package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
)

func displayMenu(app *app.App, chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, "Главное меню")
	msg.ReplyMarkup = KeyboardMain()
	return app.Bot.Reply(msg)
}
