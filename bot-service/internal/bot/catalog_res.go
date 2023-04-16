package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	c "github.com/serhiq/skye-trading-bot/internal/bot/commands"
	_type "github.com/serhiq/skye-trading-bot/pkg/type"
	"github.com/serhiq/skye-trading-bot/pkg/type/product"
)

const SELECT_CATEGORY_MESSAGE = "Выберите категорию"

const CLICK_ON_FOLDER = "@"
const CLICK_ON_PRODUCT_ITEM = "#"
const CLICK_ON_BACK = "s@f"

func ClickOnFolder(uuid string) *c.UserCommand {
	return &c.UserCommand{
		Command: CLICK_ON_FOLDER,
		Uuid:    uuid,
	}
}
func ClickOnBackInFolder(uuid string) *c.UserCommand {
	return &c.UserCommand{
		Command: CLICK_ON_BACK,
		Uuid:    uuid,
	}
}

func ClickOnProductItem(uuid string) *c.UserCommand {
	return &c.UserCommand{
		Command: CLICK_ON_PRODUCT_ITEM,
		Uuid:    uuid,
	}
}

func Keyboard(menuitems []*product.Product, isRoot bool) (tgbotapi.InlineKeyboardMarkup, error) {
	buttons := []tgbotapi.InlineKeyboardButton{}

	for _, menuitem := range menuitems {
		var title, command string
		var err error
		if menuitem.Group {
			command, err = ClickOnFolder(menuitem.UUID).ToJson()
			title = "🗀  " + menuitem.Name
		} else {
			command, err = ClickOnProductItem(menuitem.UUID).ToJson()
			title = menuitem.Name + " - " + _type.FormatPriceWithCurrency(menuitem.Price)
		}

		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, err
		}

		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(title, command))
	}

	rows := chunkSlice(buttons, calculateColums(menuitems))

	command, err := ClickOnFolder("").ToJson()
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	if !isRoot {
		backButtons := []tgbotapi.InlineKeyboardButton{}
		backButtons = append(backButtons, tgbotapi.NewInlineKeyboardButtonData("<< Назад", command))

		rows = append(rows, backButtons)
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		rows...,
	), nil
}

func calculateColums(menuitems []*product.Product) int {
	for _, menuitem := range menuitems {
		if !menuitem.Group {
			return 1
		}
	}
	return 2

}

func chunkSlice(items []tgbotapi.InlineKeyboardButton, chunkSize int) (chunks [][]tgbotapi.InlineKeyboardButton) {
	for chunkSize < len(items) {
		chunks = append(chunks, items[0:chunkSize])
		items = items[chunkSize:]
	}
	return append(chunks, items)
}
