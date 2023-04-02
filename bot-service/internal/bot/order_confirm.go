package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
)

func DisplayOrderConfirm(app *app.App, chatId int64) error {
	session, err := app.RepoChat.GetOrCreateChat(chatId)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)

	}

	order := session.GetDraftOrder()
	if order.IsEmpty() {
		msg := tgbotapi.NewMessage(chatId, EMPTY_CART_MESSAGE)
		return app.Reply(msg)
	}

	msg := tgbotapi.NewMessage(session.ChatId, ASK_ORDER_CONFIRM_MESSAGE+FormatDisplayConfirm(session, order).String())
	//msg := tgbotapi.NewMessage(session.ChatId, ms)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = MakeKeyboardConfirmOrder()

	return app.Reply(msg)
}

func ClickOnConfirm(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var deleteMsg = tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	err := app.Reply(deleteMsg)
	if err != nil {
		fmt.Printf("error delete message: %s", err)
	}

	// todo send order hereme

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)
	}

	session.NewOrder()
	err = app.RepoChat.UpdateChat(session)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, ORDER_CONFIRM_MESSAGE_TITLE+ORDER_CONFIRM_MESSAGE_BODY)
	msg.ReplyMarkup = KeyboardMain()
	return app.Reply(msg)
}

func ClickOnCancel(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var deleteMsg = tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	err := app.Reply(deleteMsg)
	if err != nil {
		fmt.Printf("error delete message: %s", err)
	}

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)
	}

	session.NewOrder()
	err = app.RepoChat.UpdateChat(session)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Заказ отменен")
	msg.ReplyMarkup = KeyboardMain()
	return app.Reply(msg)
}
