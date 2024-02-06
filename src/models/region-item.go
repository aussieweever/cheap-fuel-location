package models

type RegionItem struct {
	Region string      `json:"region"`
	Prices []PriceItem `json:"prices"`
}
