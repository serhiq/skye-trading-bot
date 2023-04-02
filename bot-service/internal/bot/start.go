package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
)

/*
   /start
*/

func StartCommand(app *app.App, message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, START_MESSAGE)
	session, err := app.RepoChat.GetOrCreateChat(message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)
	}

	if !session.HaveUserPhone() {
		err := app.Reply(msg)
		if err != nil {
			return err
		}

		requestContact := tgbotapi.NewMessage(message.Chat.ID, REQUEST_CONTACT_PHONE_MESSAGE)
		requestContact.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact(SEND_PHONE_BUTTON),
		))

		return app.Reply(requestContact)
	} else {
		msg.ReplyMarkup = KeyboardMain()
		err = app.Reply(msg)
		if err != nil {
			return err
		}

		return DisplayMenuHandler(app, message)
	}
}
