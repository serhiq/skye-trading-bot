package contorller

import (
	domainOrder "github.com/serhiq/skye-trading-bot/pkg/type/order"
	domainProduct "github.com/serhiq/skye-trading-bot/pkg/type/product"
)

type OrderSender interface {
	Send(order *domainOrder.Order) (externalUuid string, err error)
}

//type OrderStatusChecker interface {
//	CheckStatus(order *Order) (string, error)
//}

//type OrderHistoryRepository interface {
//	GetLast(count int) []*domainOrder.Order
//}

type OrderRepository interface {
	Insert(order *domainOrder.Order) error
	//Get(id int64) (*domainOrder.Order, error)
	//Update(order *domainOrder.Order) error
	//Delete(id int64) error
}

type OrderController interface {
	SendOrder(order *domainOrder.Order) (number string, err error)
}

type ProductController interface {
	GetProductByParent(parentUuid string) ([]*domainProduct.Product, error)
	GetProductByUuid(uuid string) (*domainProduct.Product, error)
	StartSync()
	StopSync()
}

type ProductRepository interface {
	Import(products []*domainProduct.Product) error
	GetProductByParent(parentUuid string) ([]*domainProduct.Product, error)
	GetProductByUuid(uuid string) (*domainProduct.Product, error)
}
