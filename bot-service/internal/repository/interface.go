package repository

import (
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
)

type ChatRepository interface {
	InsertChat(chat *chat.Chat) error
	UpdateChat(chat *chat.Chat) error
	GetChat(id int64) (*chat.Chat, error)
	GetOrCreateChat(id int64) (*chat.Chat, error)
	DeleteChat(id string) error
}
