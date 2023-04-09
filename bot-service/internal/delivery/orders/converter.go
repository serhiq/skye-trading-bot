package orderClient

import (
	"github.com/serhiq/skye-trading-bot/internal/bot"
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
	"strconv"
	"strings"
	"time"
)

func ToResponse(inputOrder *domainOrder.Order) *OrderRequest {

	positions := []OrderPosition{}

	for _, p := range inputOrder.Positions {
		positions = append(positions, OrderPosition{
			ProductUUID:       p.ProductUUID,
			Name:              p.ProductName,
			Price:             p.Price,
			PriceWithDiscount: p.Price,
			Quantity:          p.Quantity,
		})
	}

	commentBuilder := strings.Builder{}
	commentBuilder.WriteString("Способ доставки: ")
	commentBuilder.WriteString(bot.FormantDescription(inputOrder.Details.DeliveryOptions))
	commentBuilder.WriteString("\n")
	commentBuilder.WriteString("Время: ")
	commentBuilder.WriteString(bot.FormantDescription(inputOrder.Details.DeliveryTime))
	commentBuilder.WriteString("\n")

	if inputOrder.Details.DeliveryLocation != "" {
		commentBuilder.WriteString("Адрес: ")
		commentBuilder.WriteString(bot.FormantDescription(inputOrder.Details.DeliveryLocation))
		commentBuilder.WriteString("\n")
	}

	if inputOrder.Details.UserComment != "" {
		commentBuilder.WriteString("\n")
		commentBuilder.WriteString("Комментарий: \n")
		commentBuilder.WriteString(inputOrder.Details.UserComment)
		commentBuilder.WriteString("\n")
	}

	order := OrderRequest{
		UUID:        "",
		Type:        "SELL",
		Number:      inputOrder.Number,
		Period:      time.Now().Unix(),
		State:       "new",
		Client:      inputOrder.Contacts.Name,
		ClientPhone: inputOrder.Contacts.Phone,
		Comment:     commentBuilder.String(),
		ID:          strconv.FormatInt(time.Now().UnixMilli(), 32),
		Positions:   positions,
		Delivery:    Delivery{},
	}

	return &order
}
