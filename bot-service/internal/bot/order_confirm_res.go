package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
	"github.com/serhiq/skye-trading-bot/pkg/type/order"
	"strings"
)

const (
	ASK_ORDER_CONFIRM_MESSAGE = "Подтвердите информацию о заказе: \n\n"
	BUTTON_CONFIRM_ORDER      = "✅ Подтвердить заказ"
	BUTTON_CANCEL_ORDER       = "Отменить заказ"
)

const (
	CONFIRM_COMMAND = "yes"
	CANCEL_COMMAND  = "no"
)

const (
	ORDER_CONFIRM_MESSAGE_BODY = "\nВ ближайшее время наш оператор позвонит Вам для подтверждения заказа 📞"
)

func FormatConfirmMessage(number string) string {
	sb := strings.Builder{}
	sb.WriteString("✅ Заказ №")
	sb.WriteString(number)
	sb.WriteString(" подтвержден! \n")
	sb.WriteString(ORDER_CONFIRM_MESSAGE_BODY)
	return sb.String()
}

func SetConfirmOrder() *commands.UserCommand {
	return &commands.UserCommand{
		Command: CONFIRM_COMMAND,
	}
}

func SetCancelOrder() *commands.UserCommand {
	return &commands.UserCommand{
		Command: CANCEL_COMMAND,
	}
}

func MakeKeyboardConfirmOrder() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BUTTON_CONFIRM_ORDER, SetConfirmOrder().ToJson())),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BUTTON_CANCEL_ORDER, SetCancelOrder().ToJson())),
	)
}

func FormatDisplayConfirm(session *chat.Chat, order *order.Order) *strings.Builder {
	headerBuilder := strings.Builder{}
	headerBuilder.WriteString("<b>Имя: ")
	headerBuilder.WriteString(session.NameUser)
	headerBuilder.WriteString("\nТелефон: ")
	headerBuilder.WriteString(session.PhoneUser)
	headerBuilder.WriteString("</b>\n")
	headerBuilder.WriteString(order.OrderDescriptionNew())
	headerBuilder.WriteString("\n")
	headerBuilder.WriteString("Способ доставки: ")
	headerBuilder.WriteString(FormantDescription(order.Details.DeliveryOptions))
	headerBuilder.WriteString("\n")
	headerBuilder.WriteString("Время: ")
	headerBuilder.WriteString(FormantDescription(order.Details.DeliveryTime))
	headerBuilder.WriteString("\n")

	if order.Details.DeliveryLocation != "" {
		headerBuilder.WriteString("Адрес: ")
		headerBuilder.WriteString(order.Details.DeliveryLocation)
		headerBuilder.WriteString("\n")
	}

	if order.Details.UserComment != "" {
		headerBuilder.WriteString("\n")
		headerBuilder.WriteString("Комментарий: \n")
		headerBuilder.WriteString(order.Details.UserComment)
		headerBuilder.WriteString("\n")
	}

	return &headerBuilder
}
