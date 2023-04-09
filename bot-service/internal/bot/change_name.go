package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
)

func ClickOnChangeName(app *app.App, callback *tgbotapi.CallbackQuery) error {

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)
	}

	session.ChatState = chat.STATE_CHANGE_NAME
	err = app.RepoChat.UpdateChat(session)
	if err != nil {
		return err
	}

	requestContact := tgbotapi.NewMessage(callback.Message.Chat.ID, REQUEST_EDIT_NAME_MESSAGE)
	requestContact.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return app.Bot.Reply(requestContact)
}

const (
	REQUEST_EDIT_NAME_MESSAGE = "Введите имя"
)
