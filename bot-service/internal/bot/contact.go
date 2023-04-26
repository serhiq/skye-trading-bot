package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
)

func ContactHandler(app *app.App, message *tgbotapi.Message) error {
	session, err := app.RepoChat.GetOrCreateChat(message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err.Error())

	}

	var phone = message.Contact.PhoneNumber

	if phone != "" {
		if session.ChatState == chat.STATE_CHANGE_PHONE {
			return updatePhoneAndDisplayProfile(app, phone, session)
		} else {
			return SavePhoneAndInputName(app, phone, session)
		}
	} else {
		return nil

	}
}

func ContactInputHandler(app *app.App, input string, session *chat.Chat) error {
	ok, errMsg := chat.ValidateRussianPhoneNumber(input)

	if ok {
		if session.ChatState == chat.STATE_CHANGE_PHONE {
			return updatePhoneAndDisplayProfile(app, input, session)
		} else {
			return SavePhoneAndInputName(app, input, session)
		}

	} else {
		msg := tgbotapi.NewMessage(session.ChatId, errMsg)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		return app.Bot.Reply(msg)
	}
}

func SavePhoneAndInputName(app *app.App, input string, session *chat.Chat) error {
	session.PhoneUser = input
	session.ChatState = chat.STATE_INPUT_NAME
	err := app.RepoChat.UpdateChat(session)
	if err != nil {
		var errorMsg = tgbotapi.NewMessage(session.ChatId, ADD_ACCOUNT_MESSAGE_ERROR)
		return app.Bot.Reply(errorMsg)
	}
	msg := tgbotapi.NewMessage(session.ChatId, "Укажите Ваше имя, пожалуйста")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	return app.Bot.Reply(msg)
}

func updatePhoneAndDisplayProfile(app *app.App, input string, session *chat.Chat) error {
	session.PhoneUser = input
	session.ChatState = chat.STATE_PREPARE_ORDER
	err := app.RepoChat.UpdateChat(session)
	if err != nil {
		var errorMsg = tgbotapi.NewMessage(session.ChatId, ADD_ACCOUNT_MESSAGE_ERROR)
		return app.Bot.Reply(errorMsg)
	}
	return displayProfile(app, session.ChatId, session)
}

func NameInputHandler(app *app.App, input string, session *chat.Chat) error {
	if input != "" {
		var inputState = session.ChatState
		session.NameUser = input
		session.ChatState = chat.STATE_PREPARE_ORDER
		err := app.RepoChat.UpdateChat(session)
		if err != nil {
			return err
		}

		if inputState == chat.STATE_CHANGE_NAME {
			return displayProfile(app, session.ChatId, session)
		} else {
			return SayWelcome(app, input, session)
		}
	} else {
		msg := tgbotapi.NewMessage(session.ChatId, "Имя не может быть пустым")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		return app.Bot.Reply(msg)
	}
}

func SayWelcome(app *app.App, input string, session *chat.Chat) error {
	msg := tgbotapi.NewMessage(session.ChatId, FormatHelloMessage(input))
	msg.ReplyMarkup = KeyboardMain()
	return app.Bot.Reply(msg)
}
