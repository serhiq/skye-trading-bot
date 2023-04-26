package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	"github.com/serhiq/skye-trading-bot/internal/config"
	"github.com/serhiq/skye-trading-bot/pkg/type/product"
	"io/ioutil"
	"path/filepath"
)

func DisplayMenuHandler(app *app.App, message *tgbotapi.Message) error {

	menuMsg := tgbotapi.NewMessage(message.Chat.ID, SELECT_CATEGORY_MESSAGE)

	items, err := app.ProductController.GetProductByParent("")
	if err != nil {
		return fmt.Errorf("Failed to get menu  %s", err.Error())
	}
	keyboard, err := Keyboard(items, true)
	if err != nil {
		return err
	}

	menuMsg.ReplyMarkup = keyboard
	return app.Bot.Reply(menuMsg)
}

///////////////////////////////////////////////////////////

func ClickOnItemCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var c = commands.New(callback.Data)
	menuItem, err := app.ProductController.GetProductByUuid(c.Uuid)
	if err != nil {
		return fmt.Errorf("Failed to get menu item %s, err:  %s", c.Uuid, err.Error())
	}

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err.Error())
	}
	count := session.GetDraftOrder().CounterPosition(menuItem.UUID)

	text := FormatMenuItem(menuItem, count)

	keyboard, err := MakePositionKeyboard(menuItem)
	if err != nil {
		return err
	}

	if menuItem.UUID != "" {
		src := filepath.Join(config.PreviewCachePatch, menuItem.UUID)

		file, err := ioutil.ReadFile(src)
		if err != nil {
			fmt.Printf("bot: error loading image %s", err.Error())
		} else {

			var deleteMsg = tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
			err := app.Bot.Reply(deleteMsg)
			if err != nil {
				fmt.Printf("error delete message: %s", err.Error())
			}

			photoFileBytes := tgbotapi.FileBytes{
				Name:  "picture",
				Bytes: file,
			}

			photoConfig := tgbotapi.NewPhoto(session.ChatId, photoFileBytes)
			photoConfig.Caption = text
			photoConfig.ReplyMarkup = keyboard
			photoConfig.ParseMode = tgbotapi.ModeHTML

			return app.Bot.Reply(photoConfig)
		}
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(session.ChatId, callback.Message.MessageID, text, keyboard)
	msg.ParseMode = tgbotapi.ModeHTML
	return app.Bot.Reply(msg)
}

func MakePositionKeyboard(menuItem *product.Product) (tgbotapi.InlineKeyboardMarkup, error) {
	addCommand, err := AddPosition(menuItem.UUID).ToJson()
	decreaseCommand, err := DecreasePosition(menuItem.UUID).ToJson()
	backCommand, err := ClickOnBackInFolder(menuItem.ParentUUID).ToJson()

	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(formatMenuItemBtn(INCREASE_POSITION_BUTTON, menuItem), addCommand),
			tgbotapi.NewInlineKeyboardButtonData(formatMenuItemBtn(DECREASE_POSITION_BUTTON, menuItem), decreaseCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<< Назад", backCommand),
		)), nil
}

//////////////////////////////////////////////////////////////////////////////////

func ClickOnFolderCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {
	var c = commands.New(callback.Data)
	var items, err = app.ProductController.GetProductByParent(c.Uuid)
	if err != nil {
		return fmt.Errorf("Failed to get menu  %s", err.Error())
	}

	keyboard, err := Keyboard(items, c.Uuid == "")
	if err != nil {
		return err
	}
	return app.Bot.Reply(tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, keyboard))
}

///////////////////////////////////

func ClickOnBackInFolderCallbackHandler(app *app.App, callback *tgbotapi.CallbackQuery) error {

	var deleteMsg = tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	err := app.Bot.Reply(deleteMsg)
	if err != nil {
		fmt.Printf("error delete message: %s", err.Error())
	}

	var c = commands.New(callback.Data)

	items, err := app.ProductController.GetProductByParent(c.Uuid)
	if err != nil {
		return fmt.Errorf("Failed to get menu  %s", err.Error())
	}

	keyboard, err := Keyboard(items, c.Uuid == "")
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, SELECT_CATEGORY_MESSAGE)
	msg.ReplyMarkup = keyboard
	return app.Bot.Reply(msg)
}
