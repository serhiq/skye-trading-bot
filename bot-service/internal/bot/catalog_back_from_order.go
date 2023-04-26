package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	_type "github.com/serhiq/skye-trading-bot/pkg/type"
)

func BackToCatalogFromOrder(app *app.App, message *tgbotapi.Message) error {
	session, err := app.RepoChat.GetOrCreateChat(message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err.Error())

	}

	order := session.GetDraftOrder()

	msg := tgbotapi.NewMessage(message.Chat.ID, "Возврат к меню")
	msg.ReplyMarkup = MakeOrderKeyboard(_type.FormatPriceWithCurrency(order.CalculateTotal()))

	err = app.Bot.Reply(msg)
	if err != nil {
		return err
	}

	return DisplayMenuHandler(app, message)
}
