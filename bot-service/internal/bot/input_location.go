package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
)

func InputLocationHandler(app *app.App, input string, session *chat.Chat) error {

	if input != "" {
		order := session.GetDraftOrder()
		order.Details.DeliveryLocation = input

		session.ChatState = chat.INPUT_DELIVERY_TIME
		strOrder, err := order.ToJson()
		if err != nil {
			return fmt.Errorf("json error for order  =%s", err.Error())
		}
		session.OrderStr = strOrder
		err = app.RepoChat.UpdateChat(session)
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(session.ChatId, TIME_QUESTION)

		keyboard, err := KeyboardDeliveryTime()
		if err != nil {
			return err
		}

		msg.ReplyMarkup = keyboard
		return app.Bot.Reply(msg)

	} else {
		msg := tgbotapi.NewMessage(session.ChatId, "Адрес доставки не может быть пустым")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		return app.Bot.Reply(msg)
	}
}
