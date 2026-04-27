package domain

import (
	"time"
)

type GetNewestSalesOut struct {
	ItemName string
	Price    float64
	Wear     float64
	DateSold time.Time
}

type ItemWearSource struct {
	WearID int
	Name   string
	Wear   string
}

type ItemWearSale struct {
	WearID   int
	Price    float64
	Wear     float64
	DateSold time.Time
}
