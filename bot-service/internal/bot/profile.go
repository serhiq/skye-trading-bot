package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
)

func ProfileMenuHandler(app *app.App, message *tgbotapi.Message) error {
	session, err := app.RepoChat.GetOrCreateChat(message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)

	}

	return displayProfile(app, message.Chat.ID, session)
}

func displayProfile(app *app.App, chatId int64, session *chat.Chat) error {
	msg := tgbotapi.NewMessage(chatId, FormatProfileMessage(session).String())
	msg.ReplyMarkup = MakeKeyboardProfileOrder()
	msg.ParseMode = tgbotapi.ModeHTML
	return app.Reply(msg)
}
