package restoranClient

import (
	"github.com/serhiq/skye-trading-bot/internal/bot"
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
	evo "github.com/softc24/evotor-resto-go"
	"strconv"
	"strings"
	"time"
)

func ToResponse(inputOrder *domainOrder.Order) *evo.Order {

	positions := []evo.OrderPosition{}

	for _, p := range inputOrder.Positions {
		positions = append(positions, evo.OrderPosition{
			ProductUUID:       p.ProductUUID,
			ProductName:       p.ProductName,
			Price:             evo.Money(p.Price),
			PriceWithDiscount: evo.Money(p.Price),
			Quantity:          evo.Quantity(p.Quantity),
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

	order := evo.MakeOrder(strconv.FormatInt(time.Now().UnixMilli(), 32), commentBuilder.String(), evo.Contacts{
		Phone: inputOrder.Contacts.Phone,
	}, positions)

	return &order
}
