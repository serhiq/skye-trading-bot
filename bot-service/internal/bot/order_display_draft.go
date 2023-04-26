package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	_type "github.com/serhiq/skye-trading-bot/pkg/type"
)

func DisplayOrderHandler(app *app.App, message *tgbotapi.Message) error {
	return displayOrderWithMenu(app, message.Chat.ID)
}

func displayOrderWithMenu(app *app.App, chatId int64) error {
	session, err := app.RepoChat.GetOrCreateChat(chatId)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err.Error())

	}

	order := session.GetDraftOrder()
	if order.IsEmpty() {
		msg := tgbotapi.NewMessage(chatId, EMPTY_CART_MESSAGE)
		return app.Bot.Reply(msg)
	}

	msg := tgbotapi.NewMessage(session.ChatId, FormatDisplayDraft(session, order).String())
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = KeyboardOrder()

	err = app.Bot.Reply(msg)
	if err != nil {
		return err
	}
	text, keyboard, err := DisplayEditOrder(app, chatId)
	if err != nil {
		return err
	}
	// позиций в заказе нет
	if keyboard == nil {
		msg := tgbotapi.NewMessage(session.ChatId, text)
		return app.Bot.Reply(msg)
	} else {
		msg := tgbotapi.NewMessage(session.ChatId, text)
		msg.ReplyMarkup = keyboard
		return app.Bot.Reply(msg)
	}
}

//////////////////////////////////////////////////////////////////////////////////
func ClickOnEditPositionCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var c = commands.New(callback.Data)
	return displayPositionEditMenu(app, c.Uuid, callback.Message.Chat.ID, callback.Message.MessageID)
}

func displayPositionEditMenu(app *app.App, productUuid string, chatId int64, messageId int) error {
	menuItem, err := app.ProductController.GetProductByUuid(productUuid)
	if err != nil {
		return fmt.Errorf("Failed to get menu item %s, err:  %s", productUuid, err.Error())
	}

	session, err := app.RepoChat.GetOrCreateChat(chatId)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err.Error())
	}

	order := session.GetDraftOrder()
	var position = order.FindPosition(productUuid)
	if position == nil {
		return fmt.Errorf("Failed to get position  %s", productUuid)
	}

	var totalPosition = position.PriceWithDiscount * uint64(position.Quantity)

	//4шт. x 550 = 2200 ₽
	var title = fmt.Sprintf("%s\n%d x %s = %s", position.ProductName, position.Quantity, _type.FormatPrice(position.PriceWithDiscount), _type.FormatPriceWithCurrency(totalPosition))

	keyboard, err := MakePositionEditKeyboard(menuItem.UUID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(session.ChatId, messageId, title, keyboard)
	return app.Bot.Reply(msg)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
func ClickDisplayEditOrder(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var deleteMsg = tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	err := app.Bot.Reply(deleteMsg)
	if err != nil {
		fmt.Printf("error delete message: %s", err)
	}
	answer := tgbotapi.CallbackConfig{
		CallbackQueryID: callback.ID,
	}
	err = app.Bot.Reply(answer)
	if err != nil {
		fmt.Printf("callback error: %s", err)
	}

	return displayOrderWithMenu(app, callback.Message.Chat.ID)
}

func DisplayEditOrder(app *app.App, ChatId int64) (string, *tgbotapi.InlineKeyboardMarkup, error) {

	session, err := app.RepoChat.GetOrCreateChat(ChatId)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to get chat  %s", err.Error())

	}

	order := session.GetDraftOrder()

	if order.IsEmpty() {
		//msg := tgbotapi.NewMessage(ChatId, "Корзина пуста")
		return EMPTY_CART_MESSAGE, nil, nil

		//return app.Reply(msg)
	}

	//posMsg := tgbotapi.NewMessage(session.ChatId, TEXT_EDIT_QUANTITY_MESSAGE)
	editKeyboard, err := MakeEditOrderKeyboard(order)
	if err != nil {
		return "", nil, err
	}

	return TEXT_EDIT_QUANTITY_MESSAGE, &editKeyboard, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

func ClickOnIncreasePositionEditOrderCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {
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

	err = app.Bot.Reply(msg)
	if err != nil {
		return err
	}

	return displayPositionEditMenu(app, c.Uuid, callback.Message.Chat.ID, callback.Message.MessageID)
}
func ClickOnDecreasePositionEditOrderCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {
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
		// Send an empty response to the user
		answer := tgbotapi.CallbackConfig{
			CallbackQueryID: callback.ID,
		}
		return app.Bot.Reply(answer)
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

	err = app.Bot.Reply(msg)
	if err != nil {
		return err
	}

	if resultQuantity == 0 {
		return ClickDisplayEditOrder(app, callback)
	}

	return displayPositionEditMenu(app, c.Uuid, callback.Message.Chat.ID, callback.Message.MessageID)
}
