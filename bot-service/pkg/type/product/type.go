package product

import _type "github.com/serhiq/skye-trading-bot/pkg/type"

type Product struct {
	Name        string `json:"name"`
	StoreID     string `json:"storeId,omitempty"`
	UUID        string `json:"uuid"`
	ParentUUID  string `json:"parentUuid" `
	Group       bool   `json:"group"`
	Image       string `json:"image"`
	Description string `json:"description"`
	MeasureName string `json:"measureName,omitempty"`
	Price       uint64 // Цена в копейках
}

func (p *Product) PriceString() string {
	return _type.FormatPriceWithCurrency(p.Price)
}
