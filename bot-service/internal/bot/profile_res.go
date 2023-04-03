package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
	"strings"
)

const (
	CHANGE_PHONE_BUTTON   = "Изменить номер телефона"
	CHANGE_NAME_BUTTON    = "Изменить Имя"
	CANCEL_PROFILE_BUTTON = "<< Назад"
)

const (
	MESSAGE_PROFILE_TITLE    = "Здравствуйте, "
	MESSAGE_PROFILE_SUBTITLE = "Ваш телефон: "
	MESSAGE_PROFILE_BODY     = "Показать историю заказов?"
	BUTTON_MESSAGE_HISTORY   = "Да"
)

const (
	CHANGE_PHONE_COMMAND = "ch_phone"
	CHANGE_NAME_COMMAND  = "ch_name"
	CANCEL_FROM_PROFILE  = "cancel_from_profile"
)

func SetChangePhone() *commands.UserCommand {
	return &commands.UserCommand{
		Command: CHANGE_PHONE_COMMAND,
	}
}
func SetChangeName() *commands.UserCommand {
	return &commands.UserCommand{
		Command: CHANGE_NAME_COMMAND,
	}
}

func CancelProfile() *commands.UserCommand {
	return &commands.UserCommand{
		Command: CANCEL_FROM_PROFILE,
	}
}

func MakeKeyboardProfileOrder() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(CHANGE_PHONE_BUTTON, SetChangePhone().ToJson())),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(CHANGE_NAME_BUTTON, SetChangeName().ToJson())),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(CANCEL_PROFILE_BUTTON, CancelProfile().ToJson())),
	)
}

func FormatProfileMessage(session *chat.Chat) *strings.Builder {
	headerBuilder := strings.Builder{}
	headerBuilder.WriteString(MESSAGE_PROFILE_TITLE)
	headerBuilder.WriteString(session.NameUser)
	headerBuilder.WriteString("!\n\nВаш телефон: ")
	headerBuilder.WriteString(session.PhoneUser)
	headerBuilder.WriteString("\n")
	return &headerBuilder
}
