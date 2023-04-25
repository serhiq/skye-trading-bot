package order

import (
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	Db *gorm.DB
}

func (g *Repository) Get(id string) (*domainOrder.Order, error) {
	var order Order
	err := g.Db.Preload("Positions").First(&order, "uuid = ?", id).Error
	if err != nil {
		return nil, err
	}

	return mapToDomainOrder(&order), nil
}

func (g *Repository) Insert(order *domainOrder.Order) error {
	return g.Db.Save(mapToDatabaseOrder(order)).Error
}

func (g *Repository) GetLast(count int) ([]*domainOrder.Order, error) {
	var orders []*Order
	err := g.Db.Model(&Order{}).
		Preload("Positions").
		Order("updated_at DESC").
		Limit(count).
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	var result []*domainOrder.Order
	for _, order := range orders {
		result = append(result, mapToDomainOrder(order))

	}
	return result, nil
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
		Phone:       order.Contacts.Phone,
		State:       order.State,
		Comment:     order.Comment,
		Total:       total,
		Positions:   positions,
	}

}

func mapToDomainOrder(order *Order) *domainOrder.Order {
	positions := []domainOrder.Position{}

	for _, pos := range order.Positions {
		dbPos := domainOrder.Position{
			ProductUUID:       pos.ProductUUID,
			ProductName:       pos.ProductName,
			Price:             pos.Price,
			PriceWithDiscount: pos.PriceWithDiscount,
			Quantity:          int(pos.Quantity),
		}
		positions = append(positions, dbPos)
	}

	return &domainOrder.Order{
		ID:         order.UUID,
		Number:     order.Number,
		ExternalID: order.ID_EXTERNAL,
		Positions:  positions,
		State:      order.State,
		Comment:    order.Comment,
		CreatedAt:  order.Created,
		UpdatedAt:  order.Updated,
		Details:    domainOrder.Details{},
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
	ID          uint64
	UUID        string
	ID_EXTERNAL string
	Number      string
	Updated     time.Time `gorm:"column:updated_at; default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Created     time.Time `gorm:"column:created_at; default:CURRENT_TIMESTAMP()"`
	State       string
	Comment     string
	Total       uint64
	Phone       string
	Positions   []OrderPosition `gorm:"foreignKey:id_order;references:id"`
}

type OrderPosition struct {
	ID                 uint64 `gorm:"primaryKey;"`
	IdOrder            uint64 `gorm:"column:id_order;"`
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
