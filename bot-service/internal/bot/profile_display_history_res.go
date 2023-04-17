package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	"github.com/serhiq/skye-trading-bot/pkg/type/order"
	"strings"
)

const (
	REPEAT_BUTTON = "Повторить"
)

func makeHistoryOrderKeyboard(order *order.Order) (tgbotapi.InlineKeyboardMarkup, error) {
	rows := [][]tgbotapi.InlineKeyboardButton{}
	command, err := ClickOnRepeatOrder(order.ID).ToJson()
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(REPEAT_BUTTON, command)))

	return tgbotapi.NewInlineKeyboardMarkup(
		rows...,
	), nil
}

func ClickOnRepeatOrder(uuid string) *commands.UserCommand {
	return &commands.UserCommand{
		Command: CLICK_ON_REPEAT_ORDER,
		Uuid:    uuid,
	}
}

const CLICK_ON_REPEAT_ORDER = "re"

func formatDisplayHistoryOrder(order *order.Order) *strings.Builder {
	sb := strings.Builder{}
	sb.WriteString("Заказ №<b>")
	sb.WriteString(order.Number)
	sb.WriteString("</b>")
	sb.WriteString(" от ")
	sb.WriteString(order.ConvertUpdatedAtToString())
	sb.WriteString("\n")
	sb.WriteString(order.StateDescription())
	sb.WriteString("\n")
	sb.WriteString(order.OrderDescriptionNew())
	return &sb
}
