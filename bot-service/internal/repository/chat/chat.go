package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/serhiq/skye-trading-bot/pkg/type/order"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type Repository struct {
	Db *gorm.DB
}

func New(Db *gorm.DB) *Repository {
	return &Repository{
		Db: Db,
	}
}

type Chat struct {
	ChatId    int64  `gorm:"column:id;primary_key"`
	NameUser  string `gorm:"column:user_name"`
	PhoneUser string `gorm:"column:user_phone"`
	ChatState string `gorm:"column:chat_state"`
	OrderStr  string `gorm:"column:order"`
}

func NewChat(chatId int64) *Chat {
	p := new(Chat)
	p.ChatId = chatId
	p.ChatState = STATE_PREPARE_ORDER
	return p
}

func (c *Chat) HaveContact() bool {
	return len(c.NameUser) != 0 && len(c.PhoneUser) != 0
}
func (c *Chat) HaveUserName() bool {
	return len(c.NameUser) != 0
}
func (c *Chat) HaveUserPhone() bool {
	return len(c.PhoneUser) != 0
}

func (c *Chat) GetDraftOrder() *order.Order {
	var o = &order.Order{}
	err := json.Unmarshal([]byte(c.OrderStr), o)
	if err != nil {
		fmt.Print("order: unmarshal error")
	}
	return o
}

func (c *Chat) NewOrder() {
	order := &order.Order{}

	orderStr, err := order.ToJson()
	if err != nil {
		log.Print("order: marshal error")
	}
	c.OrderStr = orderStr

}

const (
	STATE_INPUT_NAME        = "STATE_INPUT_NAME"
	STATE_INPUT_PHONE       = "INPUT_PHONE"
	STATE_PREPARE_ORDER     = "PREPARE_ORDER"
	INPUT_DELIVERY_TIME     = "INPUT_DELIVERY_TIME"
	INPUT_DELIVERY_LOCATION = "INPUT_DELIVERY_LOCATION"
	INPUT_COMMENT           = "INPUT_COMMENT"
	STATE_CHANGE_NAME       = "STATE_CHANGE_NAME"
	STATE_CHANGE_PHONE      = "STATE_CHANGE_PHONE"
)

///////////////////////////////////////////////////////
type GormDatabase struct {
	Db *gorm.DB
}

func CreateGorm(db *gorm.DB) *GormDatabase {
	return &GormDatabase{Db: db}
}

func (g *Repository) InsertChat(chat *Chat) error {
	return g.Db.Create(chat).Error
}

func (r *Repository) UpdateChat(chat *Chat) error {
	return r.Db.Updates(chat).Error
}

func (r *Repository) GetChat(id int64) (*Chat, error) {
	tag := new(Chat)
	err := r.Db.Where("id = ?", id).Find(tag).Error
	return tag, err
}

func (r *Repository) GetOrCreateChat(id int64) (*Chat, error) {
	chat := new(Chat)
	err := r.Db.Where("id = ?", id).Take(&chat).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newChat := NewChat(id)
			err = r.InsertChat(newChat)
			if err != nil {
				return nil, err
			}
			return newChat, nil
		} else {
			return nil, err
		}
	}

	return chat, nil
}

func (r *Repository) DeleteChat(id string) error {
	result := r.Db.Select(clause.Associations).Unscoped().Delete(&Chat{}, id)
	return result.Error
}

func (c *Chat) IsCorrectName(name string) bool {
	return name != ""
}

func (c *Chat) IsCorrectPhone(phone string) bool {
	return ValidateRussianPhoneNumber(phone)
}
