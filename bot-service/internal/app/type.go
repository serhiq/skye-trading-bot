package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	repositoryChat "github.com/serhiq/skye-trading-bot/internal/repository"
	repositoryOrder "github.com/serhiq/skye-trading-bot/internal/repository"
	repositoryProduct "github.com/serhiq/skye-trading-bot/internal/repository"
)

type App struct {
	RepoProduct repositoryProduct.ProductRepository
	RepoChat    repositoryChat.ChatRepository
	RepoOrder   repositoryOrder.OrderRepository

	Bot *tgbotapi.BotAPI
}

func (a App) Reply(msg tgbotapi.Chattable) error {
	_, err := a.Bot.Send(msg)

	if err != nil {
		return NewErrorRespond(err)
	}

	return nil
}

/*
send an empty callback response for prevent the "waiting" icon from appearing on an inline keyboard
*/
func (a App) AnswerEmptyOnCallback(callbackID string) error {
	answer := tgbotapi.CallbackConfig{
		CallbackQueryID: callbackID,
	}
	return a.Reply(answer)

}

func (a App) ReplyWithId(msg tgbotapi.Chattable) (*tgbotapi.Message, error) {
	resultMsg, err := a.Bot.Send(msg)

	if err != nil {
		return nil, NewErrorRespond(err)
	}

	return &resultMsg, nil
}

func NewErrorRespond(err error) *ErrRespond {
	return &ErrRespond{
		err: err.Error(),
	}

}

type ErrRespond struct {
	err string
}

func (e ErrRespond) Error() string {
	return fmt.Sprintf("Failed to respond  %s", e.err)
}

func IsRespondError(err error) bool {
	_, ok := err.(ErrRespond)
	return ok
}
