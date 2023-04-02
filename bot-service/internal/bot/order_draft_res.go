package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
	"github.com/serhiq/skye-trading-bot/pkg/type/order"
	"strings"
)

const (
	CLEAR_ORDER_BUTTON = "🗑  Очистить"
	BACK_ORDER_BUTTON  = "←  Вернуться к меню"
)

const (
	EMPTY_CART_MESSAGE = "Корзина пуста"
)

func KeyboardOrder() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ORDER_BUTTON),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CLEAR_ORDER_BUTTON),
			tgbotapi.NewKeyboardButton(BACK_ORDER_BUTTON),
		),
	)
}

func FormatDisplayDraft(session *chat.Chat, order *order.Order) *strings.Builder {
	sb := strings.Builder{}
	sb.WriteString("<b>Имя: ")
	sb.WriteString(session.NameUser)
	sb.WriteString("\nТелефон: ")
	sb.WriteString(session.PhoneUser)
	sb.WriteString("</b>\n")
	sb.WriteString(order.OrderDescriptionNew())
	sb.WriteString("\nЕсли все верно выберите действие:")
	return &sb
}

const (
	TEXT_EDIT_QUANTITY_MESSAGE = "Выберите товар, чтобы изменить количество:"
)

func MakePositionEditKeyboard(productUuid string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("+", AddPositionOrder(productUuid).ToJson()),
			tgbotapi.NewInlineKeyboardButtonData("-", DecreasePositionOrder(productUuid).ToJson()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<< Назад", ClickOnBackEditOrder().ToJson()),
		))
}

func ClickOnBackEditOrder() *commands.UserCommand {
	return &commands.UserCommand{
		Command: CLICK_ON_BACK_EDIT_ORDER,
	}
}

// режим редактирования заказа

func AddPositionOrder(uuid string) *commands.UserCommand {
	return &commands.UserCommand{
		Command: CLICK_ON_ADD_POSITION,
		Uuid:    uuid,
	}
}

func DecreasePositionOrder(uuid string) *commands.UserCommand {
	return &commands.UserCommand{
		Command: CLICK_ON_DECREASE_POSITION,
		Uuid:    uuid,
	}
}

const CLICK_ON_BACK_EDIT_ORDER = "1"
const CLICK_ON_ADD_POSITION = "2"
const CLICK_ON_DECREASE_POSITION = "3"
