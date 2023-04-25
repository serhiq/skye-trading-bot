package order

import (
	"encoding/json"
	"fmt"
	_type "github.com/serhiq/skye-trading-bot/pkg/type"
	"github.com/serhiq/skye-trading-bot/pkg/type/product"
	"strconv"
	"strings"
	"time"
)

type Order struct {
	ID         string // id
	Number     string //номер, читаемвый для людей
	ExternalID string //внешний id, uuid от api
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Contacts   struct {
		Phone string
		Name  string
	}
	Positions []Position
	State     string
	Comment   string
	Details   Details
}

type Details struct {
	DeliveryOptions  string
	DeliveryLocation string
	DeliveryTime     string
	Payment          string
	UserComment      string
}

type Position struct {
	ProductUUID       string
	ProductName       string
	Price             uint64
	PriceWithDiscount uint64
	Quantity          int
}

func (c *Order) AddItem(item *product.Product) {

	for i, position := range c.Positions {
		if position.ProductUUID == item.UUID {
			c.Positions[i].Quantity = position.Quantity + 1
			return
		}
	}

	c.Positions = append(c.Positions, Position{
		ProductUUID:       item.UUID,
		ProductName:       item.Name,
		Price:             item.Price,
		PriceWithDiscount: item.Price,
		Quantity:          1,
	})
}

func (c *Order) DecreaseMenuItem(item *product.Product) int {
	for i, position := range c.Positions {
		if position.ProductUUID == item.UUID {
			if position.Quantity == 1 {
				c.Positions = append(c.Positions[:i], c.Positions[i+1:]...)
				return 0
			}

			c.Positions[i].Quantity = position.Quantity - 1
			return c.Positions[i].Quantity
		}
	}

	return -1
}

func (c *Order) CountPosition() string {
	var count uint64 = 0
	for _, position := range c.Positions {
		count = count + uint64(position.Quantity)
	}
	return strconv.FormatUint(count, 10)
}

func (c *Order) CounterPosition(uuid string) string {
	return strconv.FormatUint(c.CountItemPosition(uuid), 10)
}

func (c *Order) CountItemPosition(uuid string) uint64 {
	var count uint64 = 0
	for _, position := range c.Positions {
		if position.ProductUUID == uuid {
			count = count + uint64(position.Quantity)
		}
	}
	return count
}

func (c *Order) SumPositions() string {
	var sum uint64 = 0
	for _, position := range c.Positions {
		sum = sum + (uint64(position.Quantity) * uint64(position.PriceWithDiscount))
	}

	return strconv.FormatUint(sum, 10) + "руб"
}

func (c *Order) CalculateTotal() uint64 {
	var total uint64
	for _, position := range c.Positions {
		// Используем цену с учетом скидки, если она задана, иначе используем базовую цену
		var price uint64
		if position.PriceWithDiscount != 0 {
			price = position.PriceWithDiscount
		} else {
			price = position.Price
		}

		total += price * uint64(position.Quantity)
	}

	return total
}

func (c *Order) OrderDescription() string {
	builder := strings.Builder{}

	builder.WriteString("\n\n<b>Состав заказа:</b>")

	for i, position := range c.Positions {
		builder.WriteString("\n\n" + strconv.FormatInt(int64(i+1), 10) + ". " + position.ProductName)
		builder.WriteString("\n" + "    Кол-во: " + strconv.FormatInt(int64(position.Quantity), 10))
		builder.WriteString("\n" + "    Цена: " + strconv.FormatInt(int64(position.Price), 10) + " руб")

	}

	builder.WriteString("\n")
	builder.WriteString("\nЕсли все верно выберите действие:")

	return builder.String()
}

func (c *Order) OrderDescriptionNew() string {
	b := strings.Builder{}

	b.WriteString("\n<b>Состав заказа:</b>")

	for i, pos := range c.Positions {
		var itemPrice uint64
		if pos.PriceWithDiscount != 0 {
			itemPrice = pos.PriceWithDiscount
		} else {
			itemPrice = pos.Price
		}
		itemStr := fmt.Sprintf("\n\n%d. %s\n    Кол-во: %d\n    Цена: %s \n", i+1, pos.ProductName, pos.Quantity, _type.FormatPriceWithCurrency(itemPrice))

		b.WriteString(itemStr)
	}

	b.WriteString("\n<b>Общая сумма заказа: " + _type.FormatPriceWithCurrency(c.CalculateTotal()))
	b.WriteString("\n</b>")
	return b.String()
}

func (c *Order) ConvertUpdatedAtToString(timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)

	if err != nil {
		return "", err
	}

	localTime := c.UpdatedAt.In(loc)
	return localTime.Format("02.01.2006 15:04:05"), nil
}

func (c *Order) ToJson() (string, error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (c *Order) IsEmpty() bool {
	return len(c.Positions) == 0

}

func (c *Order) FindPosition(uuid string) *Position {
	for _, position := range c.Positions {
		if position.ProductUUID == uuid {
			return &position
		}
	}
	return nil
}

func (c *Order) StateDescription() string {
	switch c.State {
	case "new":
		return "ждет оплаты"
	case "paid":
		return "оплачен"
	case "done":
		return "выполнен"
	case "canceled":
		return "отменен"
	}
	return c.State
}
