package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	DISPLAY_MENU_BUTTON    = "🍽  Меню"
	DISPLAY_PROFILE_BUTTON = "👤  Профиль"
	DISPLAY_ORDER_BUTTON   = "🛒  Корзина "
)

const (
	CONTACT_MESSAGE_BEGIN = "Добро пожаловать, "
	CONTACT_MESSAGE_END   = "!  \n\nНажмите кнопку Меню, чтобы сделать заказ блюд."
	//CONTACT_MESSAGE_END   = "!  \n\nНажмите кнопку Меню, чтобы сделать заказ блюд. В своем профиле Вы сможете посмотреть предыдущие заказы и повторить их еще раз!"
)

const ADD_ACCOUNT_MESSAGE_ERROR = "Возникла ошибка при добавлении аккаунта"

func FormatHelloMessage(name string) string {
	return CONTACT_MESSAGE_BEGIN + name + CONTACT_MESSAGE_END
}

func KeyboardMain() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(DISPLAY_MENU_BUTTON),
			tgbotapi.NewKeyboardButton(DISPLAY_ORDER_BUTTON),
		),

		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(DISPLAY_PROFILE_BUTTON),
			tgbotapi.NewKeyboardButton(HELP_BUTTON),
		),
	)
}
