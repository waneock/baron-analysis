package domain

type ListItemSalesInput struct {
	Limit         int64
	Offset        int64
	ItemNameQuery *string
	WearName      *string
	ItemWearID    *int64
	MinPrice      *float64
	MaxPrice      *float64
	SoldFrom      *string
	SoldTo        *string
}

type ItemSale struct {
	ItemName  string
	WearName  string
	Price     float64
	WearValue float64
	SoldOn    string
}

type ListItemSalesOutput struct {
	Items  []ItemSale
	Limit  int64
	Offset int64
}

type ListItemSalesStatInput struct {
	Limit         int64
	Offset        int64
	ItemNameQuery *string
	WearName      *string
	MinPrice      *float64
	MaxPrice      *float64
	SoldFrom      *string
	SoldTo        *string
	MinSalesCount *int64
}

type ItemSalesStats struct {
	ItemID      string
	ItemName    string
	ItemWearID  int64
	WearName    string
	SalesCount  int64
	AvgPrice    float64
	MedianPrice float64
	MinPrice    float64
	MaxPrice    float64
	SoldPrices  []float64
	FirstSoldOn string
	LastSoldOn  string
}

type ListItemSalesStatOutput struct {
	Items  []ItemSalesStats
	Limit  int64
	Offset int64
}
