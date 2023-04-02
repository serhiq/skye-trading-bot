package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
)

func CreateOrderHandler(app *app.App, message *tgbotapi.Message) error {
	session, err := app.RepoChat.GetOrCreateChat(message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)

	}

	order := session.GetDraftOrder()
	if order.IsEmpty() {
		msg := tgbotapi.NewMessage(message.Chat.ID, EMPTY_CART_MESSAGE)
		return app.Reply(msg)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, DELIVERY_METHOD_QUESTION)

	msg.ReplyMarkup = MakeKeyboardDeliveryMethod()

	return app.Reply(msg)
}

func ClickOnSetDeliveryCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {
	// удаляем сообщение
	var deleteMsg = tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	err := app.Reply(deleteMsg)
	if err != nil {
		fmt.Printf("error delete message: %s", err)
	}

	var c = commands.New(callback.Data)
	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)

	}

	order := session.GetDraftOrder()
	order.Details.DeliveryOptions = c.Command

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, FormantDescription(c.Command))
	err = app.Reply(msg)
	if err != nil {
		return err
	}

	if c.Command == SELF_PICKUP_COMMAND {
		session.ChatState = chat.INPUT_DELIVERY_TIME

		strOrder, err := order.ToJson()
		if err != nil {
			return fmt.Errorf("json error for order  =%s", err)
		}
		session.OrderStr = strOrder

		err = app.RepoChat.UpdateChat(session)
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, TIME_PICKUP_QUESTION)
		msg.ReplyMarkup = KeyboardDeliveryTime()
		return app.Reply(msg)
	}

	if c.Command == DELIVERY_COMMAND {
		session.ChatState = chat.INPUT_DELIVERY_LOCATION

		strOrder, err := order.ToJson()
		if err != nil {
			return fmt.Errorf("json error for order  =%s", err)
		}
		session.OrderStr = strOrder

		err = app.RepoChat.UpdateChat(session)
		if err != nil {
			return err
		}

		msg = tgbotapi.NewMessage(callback.Message.Chat.ID, ASK_LOCATION)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		return app.Reply(msg)
	}

	return nil
}

/*
	response on callbackdata

	const (
		TIME_COMMAND_40M = "40m"
		TIME_COMMAND_120M = "120m"
		TIME_COMMAND_SOON = "soon"
	)

*/
//
//func ClickOnDeliveryTimeCallbackHandler(app  *app.App, callback *tgbotapi.CallbackQuery) error {
//	var c = commands.New(callback.Data)
//	var time = c.Command
//
//		session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
//		if err != nil {
//			return fmt.Errorf("Failed to get chat  %s", err)
//
//		}
//
//		order := session.GetDraftOrder()
//		order.Details.DeliveryTime = time
//		strOrder, err := order.ToJson()
//		if err != nil {
//			return fmt.Errorf("error command, %s, %s", c.Command, err)
//		}
//
//		session.OrderStr = strOrder
//
//		err = app.RepoChat.UpdateChat(session)
//		if err != nil {
//			return err
//		}
//
//		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, delivery.FormantDescription(c.Command))
//		err = app.Reply(msg)
//		if err != nil {
//			return err
//		}
//
//		session.ChatState = chat.INPUT_COMMENT
//		err = app.RepoChat.UpdateChat(session)
//		if err != nil {
//			return err
//		}
//
//		text := ""
//		if delivery == SELF_PICKUP_COMMAND {
//			text = DELIVERY_TIME_PICKUP_QUESTION
//		}
//
//		if delivery == DELIVERY_COMMAND {
//			text = DELIVERY_TIME_QUESTION
//
//		}
//
//		msg = tgbotapi.NewMessage(callback.Message.Chat.ID, text)
//		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
//			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(DELIVERY_BUTTON, SetTime40().ToJson())),
//			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(SELF_PICKUP_BUTTON, SetTime120().ToJson())),
//			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(SOONEST_DELIVERY_BUTTON, SetTimeSoonest().ToJson())),
//		)
//
//		return app.Reply(msg)
//	}
//
//	return nil
//}
//
