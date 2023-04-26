package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
)

func ClickOnChangePhone(app *app.App, callback *tgbotapi.CallbackQuery) error {

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err.Error())
	}

	session.ChatState = chat.STATE_CHANGE_PHONE
	err = app.RepoChat.UpdateChat(session)
	if err != nil {
		return err
	}

	requestContact := tgbotapi.NewMessage(callback.Message.Chat.ID, REQUEST_EDIT_PHONE_MESSAGE)
	requestContact.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButtonContact(SEND_PHONE_BUTTON),
	))
	return app.Bot.Reply(requestContact)
}

const (
	REQUEST_EDIT_PHONE_MESSAGE = "Введите номер телефона"
)
