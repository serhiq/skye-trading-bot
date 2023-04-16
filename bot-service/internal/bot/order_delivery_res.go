package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
)

const (
	ORDER_BUTTON             = "🚀 Оформить заказ"
	DELIVERY_METHOD_QUESTION = "Как хотите получить заказ?"
	SELF_PICKUP_BUTTON       = "Заберу сам"
	DELIVERY_BUTTON          = "Доставка"
)

const (
	DELIVERY_COMMAND    = "DELIVERY"
	SELF_PICKUP_COMMAND = "PICKUP"
)

func SetDelivery() *commands.UserCommand {
	return &commands.UserCommand{
		Command: DELIVERY_COMMAND,
	}
}

func SetPickup() *commands.UserCommand {
	return &commands.UserCommand{
		Command: SELF_PICKUP_COMMAND,
	}
}

func MakeKeyboardDeliveryMethod() (tgbotapi.InlineKeyboardMarkup, error) {
	delivery, err := SetDelivery().ToJson()
	pickup, err := SetPickup().ToJson()

	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(DELIVERY_BUTTON, delivery)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(SELF_PICKUP_BUTTON, pickup)),
	), nil
}

func FormantDescription(command string) string {
	switch command {
	case DELIVERY_COMMAND:
		return DELIVERY_BUTTON
	case SELF_PICKUP_COMMAND:
		return SELF_PICKUP_BUTTON
	case TIME_COMMAND_SOON:
		return SOONEST_DELIVERY_BUTTON
	case TIME_COMMAND_120M:
		return DELIVERY_IN_120MINS_BUTTON
	case TIME_COMMAND_40M:
		return DELIVERY_IN_120MINS_BUTTON
	default:
		return ""
	}
}
