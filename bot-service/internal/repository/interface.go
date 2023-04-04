package repository

import (
	"github.com/serhiq/skye-trading-bot/internal/repository/chat"
	product2 "github.com/serhiq/skye-trading-bot/internal/repository/product"
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
	"github.com/serhiq/skye-trading-bot/pkg/type/product"
)

type OrderRepository interface {
	Send(order *domainOrder.Order) (number string, err error)
}

type ProductRepository interface {
	ImportMenu(items []*product2.Product) error
	GetMenuItemByParent(parent string) ([]*product.Product, error)
	GetMenu(id string) (*product.Product, error)
}

type ChatRepository interface {
	InsertChat(chat *chat.Chat) error
	UpdateChat(chat *chat.Chat) error
	GetChat(id int64) (*chat.Chat, error)
	GetOrCreateChat(id int64) (*chat.Chat, error)
	DeleteChat(id string) error
}
