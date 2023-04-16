package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
)

func ClickOnSetTimeCallback(app *app.App, callback *tgbotapi.CallbackQuery) error {
	// удаляем сообщение
	var deleteMsg = tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	err := app.Bot.Reply(deleteMsg)
	if err != nil {
		fmt.Printf("error delete message: %s", err)
	}

	var c = commands.New(callback.Data)
	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)

	}

	order := session.GetDraftOrder()
	order.Details.DeliveryTime = c.Command

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, FormantDescription(c.Command))
	err = app.Bot.Reply(msg)
	if err != nil {
		return err
	}
	session.ChatState = chat.INPUT_COMMENT

	strOrder, err := order.ToJson()
	if err != nil {
		return fmt.Errorf("json error for order  =%s", err)
	}
	session.OrderStr = strOrder

	err = app.RepoChat.UpdateChat(session)
	if err != nil {
		return err
	}

	msg = tgbotapi.NewMessage(callback.Message.Chat.ID, ASK_COMMENT)

	keyboard, err := KeyboardComment()
	if err != nil {
		return err
	}

	msg.ReplyMarkup = keyboard
	return app.Bot.Reply(msg)
}
