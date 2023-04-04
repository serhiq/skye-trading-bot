package order

import (
	"github.com/serhiq/skye-trading-bot/pkg/restoClient"
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	evoClient *restoClient.RestoClient
	Db        *gorm.DB
}

func New(c *restoClient.RestoClient, Db *gorm.DB) *Repository {
	return &Repository{
		evoClient: c,
		Db:        Db,
	}
}

func (r Repository) Send(order *domainOrder.Order) (number string, err error) {
	response, err := r.evoClient.PostOrder(order)
	if err != nil {
		return "", err
	}

	result := []OrderPosition{}
	var total uint64

	for _, pos := range response.Positions {
		dbPos := OrderPosition{
			ProductUUID:       pos.ProductUUID,
			ProductName:       pos.ProductName,
			Price:             uint64(pos.Price),
			PriceWithDiscount: uint64(pos.PriceWithDiscount),
			Quantity:          uint64(pos.Quantity),
		}

		result = append(result, dbPos)
		total = +dbPos.Total()
	}

	var dbOrder = Order{
		UUID:        response.UUID,
		ID_EXTERNAL: response.ID,
		Number:      response.Number,
		CreatedAt:   response.CreatedAt.Time,
		UpdatedAt:   response.UpdatedAt.Time,
		Phone:       response.Contacts.Phone,
		State:       response.State,
		Comment:     response.Comment,
		Total:       total,
		Positions:   result,
	}

	r.Db.Save(&dbOrder)

	return response.Number, nil
}

type Order struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement:true"`
	UUID        string
	ID_EXTERNAL string
	Number      string
	CreatedAt   time.Time `json:"createdAt,omitempty"` // дата создания
	UpdatedAt   time.Time `json:"updatedAt,omitempty"` // дата обновления
	Phone       string
	State       string
	Comment     string
	Total       uint64
	Positions   []OrderPosition `gorm:"foreignKey:id"`
}

type OrderPosition struct {
	ID                 string `gorm:"primaryKey;"`
	Position           int
	IdOrder            string `gorm:"column:id_order;size:100"`
	ProductUUID        string `gorm:"column:product_uuid;"`
	ProductName        string `gorm:"column:product_name;"`
	ProductMeasureName string `gorm:"column:measure_name;"`
	Price              uint64
	PriceWithDiscount  uint64 `gorm:"column:price_with_discount;"`
	Quantity           uint64
}

func (pos OrderPosition) Total() uint64 {
	var price uint64
	if pos.PriceWithDiscount != 0 {
		price = pos.PriceWithDiscount
	} else {
		price = pos.Price
	}
	return price * pos.Quantity
}

//func GetAll(db *gorm.DB) ([]Order, error) {
//	var users []Order
//	err := db.Model(&Order{}).Preload("Positions").Find(&users).Error
//	return users, err
//}
