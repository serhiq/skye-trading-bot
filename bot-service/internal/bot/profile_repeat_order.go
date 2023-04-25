package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/app"
	"github.com/serhiq/skye-trading-bot/internal/bot/commands"
)

func ProfileRepeatOrder(app *app.App, callback *tgbotapi.CallbackQuery) error {
	err := app.Bot.AnswerEmptyOnCallback(callback.ID)
	fmt.Printf("error on send emptyCallback: %s", err)
	var c = commands.New(callback.Data)
	o, err := app.OrderController.Get(c.Uuid)
	if err != nil {
		return fmt.Errorf("Failed to get order item %s, err:  %s", c.Uuid, err)
	}

	session, err := app.RepoChat.GetOrCreateChat(callback.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("Failed to get chat  %s", err)
	}

	order := session.GetDraftOrder()

	for _, oldPosition := range o.Positions {
		product, err := app.ProductController.GetProductByUuid(oldPosition.ProductUUID)
		if err != nil {
			continue
		}

		order.AddItem(product, oldPosition.Quantity)
	}

	orderStr, err := order.ToJson()
	if err != nil {
		return err
	}
	session.OrderStr = orderStr

	err = app.RepoChat.UpdateChat(session)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Позиции добавлены в корзину")
	msg.ParseMode = tgbotapi.ModeHTML

	msg.ReplyMarkup = KeyboardMain()
	return app.Bot.Reply(msg)
}
