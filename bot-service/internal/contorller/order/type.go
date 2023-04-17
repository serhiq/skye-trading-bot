package order

import (
	"github.com/serhiq/skye-trading-bot/internal/contorller"
	"github.com/serhiq/skye-trading-bot/internal/utils"
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
	"strconv"
	"time"
)

type OrderControllerImpl struct {
	orderRepo        contorller.OrderRepository
	orderHistoryRepo contorller.OrderHistoryRepository
	orderSender      contorller.OrderSender
}

func (c *OrderControllerImpl) Get(id string) (*domainOrder.Order, error) {
	return c.orderRepo.Get(id)
}

func (c *OrderControllerImpl) GetLast(count int) ([]*domainOrder.Order, error) {
	return c.orderHistoryRepo.GetLast(count)
}

func New(orderRepo contorller.OrderRepository, orderSender contorller.OrderSender, orderHistoryRepo contorller.OrderHistoryRepository) contorller.OrderController {
	return &OrderControllerImpl{
		orderRepo:        orderRepo,
		orderSender:      orderSender,
		orderHistoryRepo: orderHistoryRepo,
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
	order.State = "new"
	err = c.orderRepo.Insert(order)
	if err != nil {
		return "", err
	}
	return orderNumber, nil
}
