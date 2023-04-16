package bot

import (
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
)

const (
	ASK_PAYMENT = "💳 Выберите метод оплаты:"

	CASH_PAYMENT_BUTTON = "💵 Наличные"
	CARD_PAYMENT_BUTTON = "💳 Безналичные"
)

const (
	COMMAND_CASH = "cash"
	COMMAND_CARD = "card"
)

func SetCardPayment() *commands.UserCommand {
	return &commands.UserCommand{
		Command: COMMAND_CARD,
	}
}

func SetCashPayment() *commands.UserCommand {
	return &commands.UserCommand{
		Command: COMMAND_CASH,
	}
}

//func KeyboardPayment() tgbotapi.InlineKeyboardMarkup {
//	return tgbotapi.NewInlineKeyboardMarkup(
//		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(CASH_PAYMENT_BUTTON, SetCashPayment().ToJson())),
//		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(CARD_PAYMENT_BUTTON, SetCardPayment().ToJson())),
//	)
//}
