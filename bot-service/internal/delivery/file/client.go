package file

import (
	"encoding/json"
	"github.com/google/uuid"
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
	"github.com/serhiq/skye-trading-bot/pkg/type/product"
	"os"
)

type FileClient struct {
}

func (f FileClient) Send(order *domainOrder.Order) (externalUuid string, err error) {
	return uuid.New().String(), nil

}

func (f FileClient) GetProducts() ([]*product.Product, error) {
	var result []*product.Product

	data, err := os.ReadFile("./data/products.json")
	if err != nil {
		return result, err
	}
	resp := []MenuItem{}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	for _, p := range resp {
		var newProduct = product.Product{
			Name:        p.Name,
			StoreID:     "",
			UUID:        p.UUID,
			ParentUUID:  p.ParentUUID,
			Group:       p.Group,
			Image:       p.ImageURL,
			Description: p.Description,
			MeasureName: p.MeasureName,
			Price:       uint64(p.Price),
		}

		result = append(result, &newProduct)
	}

	return result, nil
}

type MenuItem struct {
	UUID        string `json:"uuid"`
	Group       bool   `json:"group"`
	Name        string `json:"name"`
	Price       uint64 `json:"price,omitempty"`
	MeasureName string `json:"measureName,omitempty"`
	Description string `json:"description,omitempty"`
	ParentUUID  string `json:"parentUuid,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}
