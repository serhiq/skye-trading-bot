package evoClient

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/serhiq/skye-trading-bot/pkg/type/product"
	"log"
)

type EvotorClient struct {
	client *resty.Client

	options *Options
}

func (c EvotorClient) GetProducts() ([]*product.Product, error) {
	groups := GroupResponse{}

	endpoint := c.options.BaseUrl + c.options.Store + "/product-groups/"
	resp, err := c.client.R().
		SetHeader("Authorization", c.options.Auth).
		SetResult(&groups).
		Get(endpoint)

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		log.Println(resp)
		return nil, fmt.Errorf("client: getMenu is ERROR: code == %d, resp: == %s", resp.StatusCode(), resp)
	}

	//resultGroups := FindChildGroups(&groups.Items, c.options.MenuUuid)
	//groupMap := make(map[string]bool)
	//for _, group := range resultGroups {
	//	groupMap[group.ID] = true
	//}

	products := ProductResponse{}
	productEndpoint := c.options.BaseUrl + c.options.Store + "/products/"
	resp, err = c.client.R().
		SetHeader("Authorization", c.options.Auth).
		SetResult(&products).
		Get(productEndpoint)

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		log.Println(resp)
		return nil, fmt.Errorf("client: getMenu is ERROR: code == %d, resp: == %s", resp.StatusCode(), resp)
	}
	//resultProduct := FindAllChild(&products.Items, groupMap)

	var result []*product.Product
	for _, p := range products.Items {
		var newProduct = product.Product{
			Name:        p.Name,
			StoreID:     "",
			UUID:        p.ID,
			ParentUUID:  p.ParentID,
			Group:       false,
			Image:       "",
			MeasureName: p.MeasureName,
			Price:       uint64(p.Price * 100),
		}

		result = append(result, &newProduct)
	}

	for _, g := range groups.Items {
		var newProduct = product.Product{
			Name:        g.Name,
			StoreID:     "",
			UUID:        g.ID,
			ParentUUID:  g.ParentID,
			Group:       true,
			Image:       "",
			MeasureName: "",
			Price:       0,
		}

		result = append(result, &newProduct)
	}

	return result, nil
}

func FindAllChild(products *[]ProductItem, groupMap map[string]bool) []*ProductItem {
	var items []*ProductItem
	for _, item := range *products {
		if groupMap[item.ParentID] {
			items = append(items, &item)
		}
	}
	return items
}

type Options struct {
	Auth     string
	BaseUrl  string
	Store    string
	MenuUuid string
}

func New(client *resty.Client, options *Options) *EvotorClient {
	return &EvotorClient{
		client:  client,
		options: options,
	}
}

func FindChildGroups(groups *[]Group, parentID string) []*Group {
	var childGroups []*Group
	if parentID == "" {
		for _, group := range *groups {
			childGroups = append(childGroups, &group)
		}
		return childGroups
	}

	for _, group := range *groups {
		if group.ParentID == parentID {
			childGroups = append(childGroups, &group)
			childGroups = append(childGroups, FindChildGroups(groups, group.ID)...)
		}
	}
	return childGroups
}
