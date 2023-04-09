package order

import (
	"github.com/serhiq/skye-trading-bot/internal/contorller"
	"github.com/serhiq/skye-trading-bot/internal/utils"
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
	"strconv"
	"time"
)

type OrderControllerImpl struct {
	orderRepo contorller.OrderRepository
	//orderHistoryRepo OrderHistoryRepository
	orderSender contorller.OrderSender
}

func New(orderRepo contorller.OrderRepository, orderSender contorller.OrderSender) contorller.OrderController {
	return &OrderControllerImpl{
		orderRepo:   orderRepo,
		orderSender: orderSender,
	}
}

func (c *OrderControllerImpl) SendOrder(order *domainOrder.Order) (number string, err error) {
	order.ID = strconv.FormatInt(time.Now().UnixMilli(), 32)
	orderNumber := utils.GenerateOrderNumber()
	order.Number = orderNumber

	externalUUid, err := c.orderSender.Send(order)
	if err != nil {
		return "", err
	}

	order.ExternalID = externalUUid
	err = c.orderRepo.Insert(order)
	if err != nil {
		return "", err
	}
	return orderNumber, nil
}
