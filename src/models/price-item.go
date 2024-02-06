package models

type PriceItem struct {
	Suburb   string  `json:"suburb"`
	State    string  `json:"state"`
	Type     string  `json:"type"`
	Lng      float64 `json:"lng"`
	Price    float64 `json:"price"`
	Name     string  `json:"name"`
	Postcode string  `json:"postcode"`
	Lat      float64 `json:"lat"`
}
