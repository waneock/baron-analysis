package domain

import "time"

type ListOffersInput struct {
	Limit       int64
	Offset      int64
	AppID       *int64
	State       *int64
	NameQuery   *string
	MinPrice    *float64
	MaxPrice    *float64
	ListTime    *time.Time
	LastUpdated *time.Time
	SortBy      *string
	SortOrder   *string
}

type ListOffersOutput struct {
	Items  []Offer
	Total  int64
	Limit  int64
	Offset int64
}
