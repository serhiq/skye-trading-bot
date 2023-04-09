package order

import (
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
}

func (g *Repository) Insert(order *domainOrder.Order) error {
	return g.Db.Save(mapToDatabaseOrder(order)).Error
}

func mapToDatabaseOrder(order *domainOrder.Order) *Order {
	var total uint64
	positions := []OrderPosition{}

	for _, pos := range order.Positions {
		dbPos := OrderPosition{
			ProductUUID:       pos.ProductUUID,
			ProductName:       pos.ProductName,
			Price:             uint64(pos.Price),
			PriceWithDiscount: uint64(pos.PriceWithDiscount),
			Quantity:          uint64(pos.Quantity),
		}
		total = +dbPos.Total()
		positions = append(positions, dbPos)
	}

	return &Order{
		UUID:        order.ID,
		ID_EXTERNAL: order.ExternalID,
		Number:      order.Number,
		//CreatedAt:   order.CreatedAt.Time,
		//UpdatedAt:   order.UpdatedAt.Time,
		Phone:     order.Contacts.Phone,
		State:     order.State,
		Comment:   order.Comment,
		Total:     total,
		Positions: positions,
	}

}

func New(Db *gorm.DB) *Repository {
	return &Repository{
		Db: Db,
	}
}

type Tabler interface {
	TableName() string
}

func (Order) TableName() string {
	return "orders"
}

func (OrderPosition) TableName() string {
	return "order_positions"
}

type Order struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement:true"`
	UUID        string
	ID_EXTERNAL string `gorm:"column:id_external"`
	Number      string
	Updated     int64 `gorm:"autoUpdateTime:milli"` // Use unix milli seconds as updating time
	Created     int64 `gorm:"autoCreateTime"`       // Use unix seconds as creating time	Phone       string
	State       string
	Comment     string
	Total       uint64
	Phone       string
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
