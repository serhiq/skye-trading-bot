package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
)

const (
	ASK_COMMENT       = "Можете оставить комментарий к заказу"
	NO_COMMENT_BUTTON = "Без комментария"
)

const (
	NO_COMMENT_COMMAND = "no_comment"
)

func SetNoComment() *commands.UserCommand {
	return &commands.UserCommand{
		Command: NO_COMMENT_COMMAND,
	}
}

func KeyboardComment() (tgbotapi.InlineKeyboardMarkup, error) {
	jsonCommand, err := SetNoComment().ToJson()
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(NO_COMMENT_BUTTON, jsonCommand)),
	), nil
}
