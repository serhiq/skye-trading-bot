package restoClient

import (
	"fmt"
	"github.com/serhiq/skye-trading-bot/pkg/restoClient/order"
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
	evo "github.com/softc24/evotor-resto-go"
)

func (c RestoClient) PostOrder(o *domainOrder.Order) (*evo.Order, error) {
	request := order.ToResponse(o)

	evoResponse := &evo.Order{}

	request.State = "new"

	endpoint := c.options.BaseUrl + "/order/" + c.options.Store
	resp, err := c.client.R().
		SetHeader("Authorization", c.options.Auth).
		SetBody(request).
		SetResult(&evoResponse).
		Post(endpoint)

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("client: postOrder is ERROR: code == %d", resp.StatusCode())
	}

	// todo compare order and orderResponse

	return evoResponse, nil
}
