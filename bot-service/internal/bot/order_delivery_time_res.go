package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
)

const (
	TIME_PICKUP_QUESTION       = "Через сколько планируете забрать?\nМожете написать произвольный текст или выбрать из предложенных вариантов"
	TIME_QUESTION              = "Когда привезти заказ:"
	SOONEST_DELIVERY_BUTTON    = "Как можно быстрее"
	DELIVERY_IN_40MINS_BUTTON  = "Через 40 мин"
	DELIVERY_IN_120MINS_BUTTON = "Через 2 часа"
)

const (
	TIME_COMMAND_40M  = "40m"
	TIME_COMMAND_120M = "120m"
	TIME_COMMAND_SOON = "soon"
)

func SetTime40() *commands.UserCommand {
	return &commands.UserCommand{
		Command: TIME_COMMAND_40M,
	}
}
func SetTime120() *commands.UserCommand {
	return &commands.UserCommand{
		Command: TIME_COMMAND_120M,
	}
}
func SetTimeSoonest() *commands.UserCommand {
	return &commands.UserCommand{
		Command: TIME_COMMAND_SOON,
	}
}

func KeyboardDeliveryTime() (tgbotapi.InlineKeyboardMarkup, error) {
	soonest, err := SetTimeSoonest().ToJson()
	delivery40, err := SetTime40().ToJson()
	delivery120, err := SetTime120().ToJson()

	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(SOONEST_DELIVERY_BUTTON, soonest)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(DELIVERY_IN_40MINS_BUTTON, delivery40)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(DELIVERY_IN_120MINS_BUTTON, delivery120)),
	), nil
}
