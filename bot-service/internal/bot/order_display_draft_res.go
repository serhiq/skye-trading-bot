package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	"github.com/serhiq/skye-trading-bot/pkg/type/order"
)

func MakeEditOrderKeyboard(order *order.Order) tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{}

	for _, p := range order.Positions {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(p.ProductName, ClickOnEditPosition(p.ProductUUID).ToJson())))
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		rows...,
	)
}

func ClickOnEditPosition(uuid string) *commands.UserCommand {
	return &commands.UserCommand{
		Command: CLICK_ON_EDIT_POSITION,
		Uuid:    uuid,
	}
}

const CLICK_ON_EDIT_POSITION = "%"
