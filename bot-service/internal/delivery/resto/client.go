package restoranClient

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
	"github.com/serhiq/skye-trading-bot/pkg/type/product"
	evo "github.com/softc24/evotor-resto-go"
	"log"
)

type RestoClient struct {
	client  *resty.Client
	options *Options
}

func (c RestoClient) Send(o *domainOrder.Order) (externalUuid string, err error) {
	request := ToResponse(o)

	evoResponse := &evo.Order{}

	request.State = "new"

	endpoint := c.options.BaseUrl + "/order/" + c.options.Store
	resp, err := c.client.R().
		SetHeader("Authorization", c.options.Auth).
		SetBody(request).
		SetResult(&evoResponse).
		Post(endpoint)

	if err != nil {
		return "", err
	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("client: postOrder is ERROR: code == %d", resp.StatusCode())
	}

	return evoResponse.UUID, nil
}

func (c RestoClient) GetProducts() ([]*product.Product, error) {
	respProduct := Menu{}

	endpoint := c.options.BaseUrl + "/product/" + c.options.Store
	resp, err := c.client.R().
		SetHeader("Authorization", c.options.Auth).
		SetResult(&respProduct).
		Get(endpoint)

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		log.Println(resp)
		return nil, fmt.Errorf("client: getMenu is ERROR: code == %d, resp: == %s", resp.StatusCode(), resp)
	}

	var result []*product.Product
	for _, p := range respProduct {
		if !isProductSupported(p) {
			continue
		}

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

type Options struct {
	Auth    string
	Store   string
	BaseUrl string
}

func New(client *resty.Client, options *Options) *RestoClient {
	return &RestoClient{
		client:  client,
		options: options,
	}
}

type MenuResponse evo.MenuItem
type Menu []evo.MenuItem

func (m *MenuResponse) CanAddToOrder() bool {
	if m.Group {
		return true
	}

	if !m.Group && m.Type == "NORMAL" {
		return true
	}

	return false
}

func (m *MenuResponse) ImageNotEmpty() bool {

	if m.ImageURL != "" {
		return true
	}
	return false
}

func isProductSupported(p evo.MenuItem) bool {

	if p.IsUnavailable {
		return false
	}

	if p.Group {
		return true
	}

	if !p.Group && p.Type == "NORMAL" {
		return true
	}
	return true
}
