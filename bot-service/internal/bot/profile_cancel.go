package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
)

func ProfileOnCancel(app *app.App, callback *tgbotapi.CallbackQuery) error {
	err := app.AnswerEmptyOnCallback(callback.ID)
	fmt.Printf("error on send emptyCallback: %s", err)
	return displayMenu(app, callback.Message.Chat.ID)
}
