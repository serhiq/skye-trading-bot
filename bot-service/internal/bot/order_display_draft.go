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
		return fmt.Errorf("Failed to get chat  %s", err)

	}

	order := session.GetDraftOrder()
	if order.IsEmpty() {
		msg := tgbotapi.NewMessage(chatId, EMPTY_CART_MESSAGE)
		return app.Reply(msg)
	}

	msg := tgbotapi.NewMessage(session.ChatId, FormatDisplayDraft(session, order).String())
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = KeyboardOrder()

	err = app.Reply(msg)
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
		return app.Reply(msg)
	} else {
		msg := tgbotapi.NewMessage(session.ChatId, text)
		msg.ReplyMarkup = keyboard
		return app.Reply(msg)
	}
}

//////////////////////////////////////////////////////////////////////////////////
func ClickOnEditPositionCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var c = commands.New(callback.Data)
	return displayPositionEditMenu(app, c.Uuid, callback.Message.Chat.ID, callback.Message.MessageID)
}

func displayPositionEditMenu(app *app.App, productUuid string, chatId int64, messageId int) error {
	menuItem, err := app.RepoProduct.GetMenu(productUuid)
	if err != nil {
		return fmt.Errorf("Failed to get menu item %s, err:  %s", productUuid, err)
	}

	session, err := app.RepoChat.GetOrCreateChat(chatId)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)
	}

	order := session.GetDraftOrder()
	var position = order.FindPosition(productUuid)
	if position == nil {
		return fmt.Errorf("Failed to get position  %s", productUuid)
	}

	var totalPosition = position.PriceWithDiscount * uint64(position.Quantity)

	//4шт. x 550 = 2200 ₽
	var title = fmt.Sprintf("%s\n%d x %s = %s", position.ProductName, position.Quantity, _type.FormatPrice(position.PriceWithDiscount), _type.FormatPriceWithCurrency(totalPosition))

	msg := tgbotapi.NewEditMessageTextAndMarkup(session.ChatId, messageId, title, MakePositionEditKeyboard(menuItem.UUID))
	return app.Reply(msg)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
func ClickDisplayEditOrder(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var deleteMsg = tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	err := app.Reply(deleteMsg)
	if err != nil {
		fmt.Printf("error delete message: %s", err)
	}

	return displayOrderWithMenu(app, callback.Message.Chat.ID)
}

func DisplayEditOrder(app *app.App, ChatId int64) (string, *tgbotapi.InlineKeyboardMarkup, error) {

	session, err := app.RepoChat.GetOrCreateChat(ChatId)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to get chat  %s", err)

	}

	order := session.GetDraftOrder()

	if order.IsEmpty() {
		//msg := tgbotapi.NewMessage(ChatId, "Корзина пуста")
		return EMPTY_CART_MESSAGE, nil, nil

		//return app.Reply(msg)
	}

	posMsg := tgbotapi.NewMessage(session.ChatId, TEXT_EDIT_QUANTITY_MESSAGE)
	editKeyboard := MakeEditOrderKeyboard(order)
	posMsg.ReplyMarkup = editKeyboard

	return TEXT_EDIT_QUANTITY_MESSAGE, &editKeyboard, nil
	//return app.Reply(posMsg)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

func ClickOnIncreasePositionEditOrderCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var c = commands.New(callback.Data)
	menuItem, err := app.RepoProduct.GetMenu(c.Uuid)
	if err != nil {
		return fmt.Errorf("Failed to get menu item %s, err:  %s", c.Uuid, err)
	}

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)

	}
	order := session.GetDraftOrder()
	order.AddItem(menuItem)

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

	err = app.Reply(msg)
	if err != nil {
		return err
	}

	return displayPositionEditMenu(app, c.Uuid, callback.Message.Chat.ID, callback.Message.MessageID)
}
func ClickOnDecreasePositionEditOrderCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var c = commands.New(callback.Data)
	menuItem, err := app.RepoProduct.GetMenu(c.Uuid)
	if err != nil {
		return fmt.Errorf("Failed to get menu item %s, err:  %s", c.Uuid, err)
	}

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)

	}
	order := session.GetDraftOrder()
	order.DecreaseMenuItem(menuItem)

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

	err = app.Reply(msg)
	if err != nil {
		return err
	}

	return displayPositionEditMenu(app, c.Uuid, callback.Message.Chat.ID, callback.Message.MessageID)
}
