package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
)

func ClearOrderHandler(app *app.App, message *tgbotapi.Message) error {
	session, err := app.RepoChat.GetOrCreateChat(message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)

	}

	order := session.GetDraftOrder()
	if order.IsEmpty() {
		msg := tgbotapi.NewMessage(message.Chat.ID, EMPTY_CART_MESSAGE)
		return app.Reply(msg)
	}

	session.NewOrder()
	session.ChatState = chat.STATE_PREPARE_ORDER

	err = app.RepoChat.UpdateChat(session)
	if err != nil {
		return app.Reply(tgbotapi.NewMessage(message.Chat.ID, "Ошибка очистки заказа"))
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, CLEAR_ORDER_MESSAGE)
	msg.ReplyMarkup = MakeOrderKeyboard("0")

	err = app.Reply(msg)
	if err != nil {
		return err
	}

	return DisplayMenuHandler(app, message)
}
