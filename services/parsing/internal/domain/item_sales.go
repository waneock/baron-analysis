package domain

import "time"

type ItemSales struct {
	ItemName  string
	WearName  string
	Price     float64
	WearValue float64
	SoldOn    time.Time
}

type ListItemSalesFilter struct {
	Limit         int64
	Offset        int64
	ItemNameQuery *string
	WearName      *string
	ItemWearID    *int64
	MinPrice      *float64
	MaxPrice      *float64
	SoldFrom      *time.Time
	SoldTo        *time.Time
}

type ItemSalesStats struct {
	ItemID     string
	ItemName   string
	ItemWearID int64
	WearName   string

	SalesCount  int64
	AvgPrice    float64
	MedianPrice float64
	MinPrice    float64
	MaxPrice    float64
	SoldPrices  []float64

	FirstSoldOn time.Time
	LastSoldOn  time.Time
}

type ListItemSalesStatsFilter struct {
	Limit         int64
	Offset        int64
	ItemNameQuery *string
	WearName      *string
	MinPrice      *float64
	MaxPrice      *float64
	SoldFrom      *time.Time
	SoldTo        *time.Time
	MinSalesCount *int64
}
