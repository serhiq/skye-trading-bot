package product

import (
	domainProduct "github.com/serhiq/skye-trading-bot/pkg/type/product"
	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
}

func (g *Repository) Import(products []*domainProduct.Product) error {
	err := g.Db.Transaction(func(tx *gorm.DB) error {

		g.Db.Exec("DELETE FROM products")

		for _, item := range products {
			if err := g.Db.Create(mapToDatabaseProduct(item)).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func mapToDatabaseProduct(r *domainProduct.Product) *Product {
	return &Product{
		Name:        r.Name,
		UUID:        r.UUID,
		ParentUUID:  r.ParentUUID,
		Group:       r.Group,
		Image:       r.Image,
		MeasureName: r.MeasureName,
		Price:       r.Price,
	}

}

func (g *Repository) GetProductByParent(parentUuid string) ([]*domainProduct.Product, error) {
	var items []Product
	var products = []*domainProduct.Product{}

	err := g.Db.Where("parent_uuid = ?", parentUuid).Model(&Product{}).Find(&items).Error
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		products = append(products, ToDomain(&item))

	}

	return products, err
}

func (g *Repository) GetProductByUuid(uuid string) (*domainProduct.Product, error) {
	var item Product
	err := g.Db.Where("uuid = ?", uuid).Find(&item).Error
	return ToDomain(&item), err
}

func New(Db *gorm.DB) *Repository {
	return &Repository{
		Db: Db,
	}
}

type Product struct {
	Name        string `json:"name"`
	StoreID     string `json:"storeId,omitempty" gorm:"column:store_id"`
	UUID        string `json:"uuid"`
	ParentUUID  string `json:"parentUuid"  gorm:"column:parent_uuid"`
	Group       bool   `json:"group"`
	Image       string `json:"image"`
	MeasureName string `json:"measureName,omitempty"  gorm:"column:measure_name"`
	Description string `json:"description,omitempty"  gorm:"column:description"`
	Price       uint64 // Цена в копейках
}

type Tabler interface {
	TableName() string
}

func (Product) TableName() string {
	return "products"
}

func ToDomain(r *Product) *domainProduct.Product {
	return &domainProduct.Product{
		Name:        r.Name,
		UUID:        r.UUID,
		ParentUUID:  r.ParentUUID,
		Group:       r.Group,
		Image:       r.Image,
		MeasureName: r.MeasureName,
		Price:       r.Price,
	}
}
