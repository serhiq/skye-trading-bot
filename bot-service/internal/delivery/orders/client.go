package orderClient

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
)

type OrderClient struct {
	client  *resty.Client
	options *Options
}

func (c OrderClient) Send(o *domainOrder.Order) (externalUuid string, err error) {
	request := ToResponse(o)

	evoResponse := &OrderResponse{}

	request.State = "new"

	endpoint := c.options.BaseUrl + "/order/" + c.options.Store
	resp, err := c.client.R().
		SetHeader("Authorization", c.options.Auth).
		SetHeader("Content-Type", "application/json").
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

type Options struct {
	Auth    string
	Store   string
	BaseUrl string
}

func New(client *resty.Client, options *Options) *OrderClient {
	return &OrderClient{
		client:  client,
		options: options,
	}
}

type ProductItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	MeasureName string `json:"measure_name"`
	AllowToSell bool   `json:"allow_to_sell"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id"`
}

type OrderRequest struct {
	UUID        string          `json:"uuid"`
	Type        string          `json:"type"`
	Number      string          `json:"number"`
	Period      int64           `json:"period"`
	State       string          `json:"state"`
	Client      string          `json:"client"`
	ClientPhone string          `json:"client_phone"`
	Comment     string          `json:"comment"`
	ID          string          `json:"id"`
	Positions   []OrderPosition `json:"positions"`
	Delivery    Delivery        `json:"delivery"`
}

type OrderPosition struct {
	ProductUUID       string `json:"product_uuid"`
	Name              string `json:"name"`
	Price             uint64 `json:"price"`
	Quantity          int    `json:"quantity"`
	PriceWithDiscount uint64 `json:"priceWithDiscount"`
}

type Delivery struct {
	Date     string `json:"date"`
	Address  string `json:"address"`
	TimeFrom string `json:"timeFrom"`
	TimeTo   string `json:"timeTo"`
}

type OrderResponse struct {
	UUID string `json:"uuid"`
}
