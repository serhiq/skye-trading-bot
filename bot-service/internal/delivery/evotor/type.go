package evoClient

type ProductResponse struct {
	Items  []ProductItem `json:"items"`
	Paging struct {
	} `json:"paging"`
}

type ProductItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	MeasureName string  `json:"measure_name"`
	AllowToSell bool    `json:"allow_to_sell"`
	Description string  `json:"description"`
	ParentID    string  `json:"parent_id"`
}

type GroupResponse struct {
	Items  []Group `json:"items"`
	Paging struct {
	} `json:"paging"`
}

type Group struct {
	ID       string `json:"id"`
	ParentID string `json:"parent_id,omitempty"`
	Name     string `json:"name"`
}
