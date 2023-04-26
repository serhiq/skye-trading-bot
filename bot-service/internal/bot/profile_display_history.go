package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
)

func ProfileDisplayHistory(app *app.App, callback *tgbotapi.CallbackQuery) error {
	err := app.Bot.AnswerEmptyOnCallback(callback.ID)
	fmt.Printf("error on send emptyCallback: %s", err)

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err.Error())
	}

	orders, err := app.OrderController.GetLast(3)

	if err != nil {
		msg := tgbotapi.NewMessage(session.ChatId, SAY_CANT_GET_ORDER)
		return app.Bot.Reply(msg)
	}

	if len(orders) == 0 {
		msg := tgbotapi.NewMessage(session.ChatId, SAY_EMPTY_HISTORY)
		return app.Bot.Reply(msg)
	}

	for _, order := range orders {
		msg := tgbotapi.NewMessage(session.ChatId, formatDisplayHistoryOrder(order, app.Cfg.TimeZone).String())
		msg.ParseMode = tgbotapi.ModeHTML

		keyboard, err := makeHistoryOrderKeyboard(order)
		if err != nil {
			return err
		}
		msg.ReplyMarkup = keyboard
		err = app.Bot.Reply(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
