package models

type QueryResult struct {
	Updated int64        `json:"updated"`
	Regions []RegionItem `json:"regions"`
}
