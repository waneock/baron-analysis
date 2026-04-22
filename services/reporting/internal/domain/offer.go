package domain

import "time"

type Offer struct {
	ID          string
	Price       float64
	Commission  float64
	Name        string
	State       string
	ListTime    time.Time
	LastUpdated time.Time
	Wear        float64
	Addons      []string
}
