package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
)

func InputCommentHandler(app *app.App, input string, session *chat.Chat) error {
	if input != "" {
		order := session.GetDraftOrder()
		order.Details.UserComment = input

		strOrder, err := order.ToJson()
		if err != nil {
			return fmt.Errorf("json error for order  =%s", err.Error())
		}
		session.OrderStr = strOrder

		session.ChatState = chat.STATE_PREPARE_ORDER
		err = app.RepoChat.UpdateChat(session)
		if err != nil {
			return err
		}

		return DisplayOrderConfirm(app, session.ChatId)
	} else {
		session.ChatState = chat.STATE_PREPARE_ORDER
		err := app.RepoChat.UpdateChat(session)
		if err != nil {
			return err
		}

		// введен пустой коммент, просто идем дальше
		return DisplayOrderConfirm(app, session.ChatId)
	}

}

func ClickOnNoCommentCallback(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var deleteMsg = tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	err := app.Bot.Reply(deleteMsg)
	if err != nil {
		fmt.Printf("error delete message: %s", err)
	}

	return DisplayOrderConfirm(app, callback.Message.Chat.ID)
}
