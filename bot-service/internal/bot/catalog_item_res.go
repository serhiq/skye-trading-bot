package bot

import (
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	_type "github.com/serhiq/skye-trading-bot/pkg/type"
	"github.com/serhiq/skye-trading-bot/pkg/type/product"
	"strings"
)

//const INCREASE_POSITION_BUTTON = "+"
//const DECREASE_POSITION_BUTTON = "-"

const INCREASE_POSITION_BUTTON = "➕"
const DECREASE_POSITION_BUTTON = "➖"

func FormatMenuItem(menuItem *product.Product, count string) string {
	return menuItem.Name + "\n" + "Цена за " + menuItem.MeasureName + ": " + _type.FormatPriceWithCurrency(menuItem.Price) + "\n\n" + "В корзине: " + count
}

func formatMenuItemBtn(symbol string, menuItem *product.Product) string {
	sb := strings.Builder{}
	sb.WriteString(symbol)
	sb.WriteString(" 1 ")
	sb.WriteString(menuItem.MeasureName)
	return sb.String()
}

const COMMAND_ADD_POSITION = "+"
const COMMAND_DECREASE_POSITION = "-"

func AddPosition(uuid string) *commands.UserCommand {
	return &commands.UserCommand{
		Command: COMMAND_ADD_POSITION,
		Uuid:    uuid,
	}
}

func DecreasePosition(uuid string) *commands.UserCommand {
	return &commands.UserCommand{
		Command: COMMAND_DECREASE_POSITION,
		Uuid:    uuid,
	}
}
