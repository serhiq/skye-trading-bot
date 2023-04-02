package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
)

func ProfileOnCancel(app *app.App, callback *tgbotapi.CallbackQuery) error {
	return displayMenu(app, callback.Message.Chat.ID)
}
