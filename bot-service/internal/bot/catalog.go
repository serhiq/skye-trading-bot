package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	"github.com/serhiq/skye-trading-bot/pkg/type/product"
	"io/ioutil"
)

func DisplayMenuHandler(app *app.App, message *tgbotapi.Message) error {

	menuMsg := tgbotapi.NewMessage(message.Chat.ID, SELECT_CATEGORY_MESSAGE)

	items, err := app.RepoProduct.GetMenuItemByParent("")
	if err != nil {
		return fmt.Errorf("Failed to get menu  %s", err)
	}
	keyboard := Keyboard(items, true)
	menuMsg.ReplyMarkup = keyboard
	return app.Reply(menuMsg)
}

///////////////////////////////////////////////////////////

func ClickOnItemCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var c = commands.New(callback.Data)
	menuItem, err := app.RepoProduct.GetMenu(c.Uuid)
	if err != nil {
		return fmt.Errorf("Failed to get menu item %s, err:  %s", c.Uuid, err)
	}

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)

	}
	count := session.GetDraftOrder().CounterPosition(menuItem.UUID)

	text := FormatMenuItem(menuItem, count)
	var keyboard = MakePositionKeyboard(menuItem)

	if menuItem.UUID != "" {
		src := "./imageCache/" + menuItem.UUID

		file, err := ioutil.ReadFile(src)
		if err != nil {
			fmt.Printf("bot: error loading image %s", err)
		} else {

			var deleteMsg = tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
			err := app.Reply(deleteMsg)
			if err != nil {
				fmt.Printf("error delete message: %s", err)
			}

			photoFileBytes := tgbotapi.FileBytes{
				Name:  "picture",
				Bytes: file,
			}

			photoConfig := tgbotapi.NewPhoto(session.ChatId, photoFileBytes)
			photoConfig.Caption = text
			photoConfig.ReplyMarkup = keyboard
			photoConfig.ParseMode = tgbotapi.ModeHTML

			return app.Reply(photoConfig)
		}
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(session.ChatId, callback.Message.MessageID, text, keyboard)
	msg.ParseMode = tgbotapi.ModeHTML
	return app.Reply(msg)
}

func MakePositionKeyboard(menuItem *product.Product) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(formatMenuItemBtn(INCREASE_POSITION_BUTTON, menuItem), AddPosition(menuItem.UUID).ToJson()),
			tgbotapi.NewInlineKeyboardButtonData(formatMenuItemBtn(DECREASE_POSITION_BUTTON, menuItem), DecreasePosition(menuItem.UUID).ToJson()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<< Назад", ClickOnBackInFolder(menuItem.ParentUUID).ToJson()),
		))
}

//////////////////////////////////////////////////////////////////////////////////

func ClickOnFolderCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var c = commands.New(callback.Data)
	var items, err = app.RepoProduct.GetMenuItemByParent(c.Uuid)
	if err != nil {
		return fmt.Errorf("Failed to get menu  %s", err)
	}

	var keyboard = Keyboard(items, c.Uuid == "")
	return app.Reply(tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, keyboard))
}

///////////////////////////////////

func ClickOnBackInFolderCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {

	var deleteMsg = tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	err := app.Reply(deleteMsg)
	if err != nil {
		fmt.Printf("error delete message: %s", err)
	}

	var c = commands.New(callback.Data)

	items, err := app.RepoProduct.GetMenuItemByParent(c.Uuid)
	if err != nil {
		return fmt.Errorf("Failed to get menu  %s", err)
	}

	var keyboard = Keyboard(items, c.Uuid == "")

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, SELECT_CATEGORY_MESSAGE)
	msg.ReplyMarkup = keyboard
	return app.Reply(msg)
}
