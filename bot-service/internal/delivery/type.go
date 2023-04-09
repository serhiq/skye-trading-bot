package delivery

import (
	domainProduct "github.com/serhiq/skye-trading-bot/pkg/type/product"
)

type ProductProvider interface {
	GetProducts() ([]*domainProduct.Product, error)
}
