package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
)

func DisplayOrderConfirm(app *app.App, chatId int64) error {
	session, err := app.RepoChat.GetOrCreateChat(chatId)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err.Error())

	}

	order := session.GetDraftOrder()
	if order.IsEmpty() {
		msg := tgbotapi.NewMessage(chatId, EMPTY_CART_MESSAGE)
		return app.Bot.Reply(msg)
	}

	msg := tgbotapi.NewMessage(session.ChatId, ASK_ORDER_CONFIRM_MESSAGE+FormatDisplayConfirm(session, order).String())
	//msg := tgbotapi.NewMessage(session.ChatId, ms)

	msg.ParseMode = tgbotapi.ModeHTML
	keyboard, err := MakeKeyboardConfirmOrder()
	if err != nil {
		return err
	}

	msg.ReplyMarkup = keyboard

	return app.Bot.Reply(msg)
}

func ClickOnConfirm(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var deleteMsg = tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	err := app.Bot.Reply(deleteMsg)
	if err != nil {
		fmt.Printf("error delete message: %s", err)
	}

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err.Error())
	}

	order := session.GetDraftOrder()
	order.Contacts.Phone = session.PhoneUser
	order.Contacts.Name = session.NameUser

	number, err := app.OrderController.SendOrder(order)
	if err != nil {
		return err
	}

	session.NewOrder()
	err = app.RepoChat.UpdateChat(session)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, FormatConfirmMessage(number))
	msg.ParseMode = tgbotapi.ModeHTML

	msg.ReplyMarkup = KeyboardMain()
	return app.Bot.Reply(msg)
}

func ClickOnCancel(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var deleteMsg = tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	err := app.Bot.Reply(deleteMsg)
	if err != nil {
		fmt.Printf("error delete message: %s", err.Error())
	}

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err.Error())
	}

	session.NewOrder()
	err = app.RepoChat.UpdateChat(session)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Заказ отменен")
	msg.ReplyMarkup = KeyboardMain()
	return app.Bot.Reply(msg)
}
