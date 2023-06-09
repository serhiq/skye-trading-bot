package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	_type "github.com/serhiq/skye-trading-bot/pkg/type"
)

func ClickOnDecreasePositionCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var c = commands.New(callback.Data)
	menuItem, err := app.ProductController.GetProductByUuid(c.Uuid)
	if err != nil {
		return fmt.Errorf("Failed to get menu item %s, err:  %s", c.Uuid, err.Error())
	}

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err.Error())
	}

	order := session.GetDraftOrder()

	resultQuantity := order.DecreaseMenuItem(menuItem)
	if resultQuantity == -1 {
		return app.Bot.AnswerEmptyOnCallback(callback.ID)
	}

	var msgText = "удалена позиция " + menuItem.Name + " " + menuItem.PriceString()
	strOrder, err := order.ToJson()
	if err != nil {
		return fmt.Errorf("Decrease position command, json error for product  =%s", c.Uuid)
	}

	session.OrderStr = strOrder

	err = app.RepoChat.UpdateChat(session)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(session.ChatId, msgText)
	msg.ReplyMarkup = MakeOrderKeyboard(_type.FormatPriceWithCurrency(order.CalculateTotal()))
	err = app.Bot.Reply(msg)
	if err != nil {
		return err
	}

	count := order.CounterPosition(menuItem.UUID)

	text := FormatMenuItem(menuItem, count)
	keyboard, err := MakePositionKeyboard(menuItem)

	if callback.Message.Caption != "" {
		editor := tgbotapi.NewEditMessageCaption(session.ChatId, callback.Message.MessageID, text)
		editor.ParseMode = tgbotapi.ModeHTML
		markup := keyboard
		editor.ReplyMarkup = &markup

		err = app.Bot.Reply(editor)
		if err != nil {
			return err
		}

	} else {
		updMsg := tgbotapi.NewEditMessageText(session.ChatId, callback.Message.MessageID, text)

		err = app.Bot.Reply(updMsg)
		if err != nil {
			return err
		}
	}
	return nil
}

func ClickOnIncreasePositionCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var c = commands.New(callback.Data)
	menuItem, err := app.ProductController.GetProductByUuid(c.Uuid)
	if err != nil {
		return fmt.Errorf("Failed to get menu item %s, err:  %s", c.Uuid, err.Error())
	}

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err.Error())

	}

	order := session.GetDraftOrder()
	order.AddItem(menuItem, 1)

	var msgText = "В заказ добавлена позиция " + menuItem.Name + " " + menuItem.PriceString()
	strOrder, err := order.ToJson()
	if err != nil {
		return fmt.Errorf("Decrease position command, json error for product  =%s", c.Uuid)
	}

	session.OrderStr = strOrder

	err = app.RepoChat.UpdateChat(session)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(session.ChatId, msgText)
	msg.ReplyMarkup = MakeOrderKeyboard(_type.FormatPriceWithCurrency(order.CalculateTotal()))
	err = app.Bot.Reply(msg)
	if err != nil {
		return err
	}

	count := order.CounterPosition(menuItem.UUID)

	text := FormatMenuItem(menuItem, count)
	keyboard, err := MakePositionKeyboard(menuItem)
	if err != nil {
		return err
	}

	if callback.Message.Caption != "" {
		editor := tgbotapi.NewEditMessageCaption(session.ChatId, callback.Message.MessageID, text)
		editor.ParseMode = tgbotapi.ModeHTML
		markup := keyboard
		editor.ReplyMarkup = &markup

		err = app.Bot.Reply(editor)
		if err != nil {
			return err
		}

	} else {
		updMsg := tgbotapi.NewEditMessageText(session.ChatId, callback.Message.MessageID, text)

		err = app.Bot.Reply(updMsg)
		if err != nil {
			return err
		}
	}

	return nil
}

func MakeOrderKeyboard(count string) interface{} {
	var textBucket = DISPLAY_ORDER_BUTTON + "(" + count + ")"

	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(DISPLAY_PROFILE_BUTTON),
			tgbotapi.NewKeyboardButton(DISPLAY_MENU_BUTTON),
		),

		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(textBucket),
		),
	)
}

//short variant keyboard
//func MakeOrderKeyboard(count string) interface{} {
//	var textBucket = DISPLAY_ORDER_BUTTON + "(" + count + ")"
//
//	return tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton(DISPLAY_MENU_BUTTON),
//		),
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton(textBucket),
//		),
//	)
//}
