package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	"github.com/serhiq/skye-trading-bot/pkg/type/order"
)

func MakeEditOrderKeyboard(order *order.Order) (tgbotapi.InlineKeyboardMarkup, error) {
	rows := [][]tgbotapi.InlineKeyboardButton{}

	for _, p := range order.Positions {
		command, err := ClickOnEditPosition(p.ProductUUID).ToJson()
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(p.ProductName, command)))
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		rows...,
	), nil
}

func ClickOnEditPosition(uuid string) *commands.UserCommand {
	return &commands.UserCommand{
		Command: CLICK_ON_EDIT_POSITION,
		Uuid:    uuid,
	}
}

const CLICK_ON_EDIT_POSITION = "%"
