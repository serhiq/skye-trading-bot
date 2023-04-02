package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
)

func AboutCommand(app *app.App, message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, HELP_MESSAGE)
	return app.Reply(msg)
}
