package order

import (
	"github.com/serhiq/skye-trading-bot/pkg/restoClient"
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
)

type Repository struct {
	evoClient *restoClient.RestoClient
}

func New(c *restoClient.RestoClient) *Repository {
	return &Repository{
		evoClient: c,
	}
}

func (r Repository) Send(order *domainOrder.Order) error {
	_, err := r.evoClient.PostOrder(order)
	if err != nil {
		return err
	}
	//	 todo check postorder and order
	return nil
}
