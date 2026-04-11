package domain

import "time"

type Offer struct {
	ID               string
	Price            int
	Commission       int
	Tax              int
	ClassID          string
	InstanceID       string
	AppID            string
	ContextID        string
	AssetID          string
	Name             string
	OfferID          string
	State            int
	EscrowEndDate    time.Time
	ListTime         time.Time
	LastUpdated      time.Time
	Wear             int
	TxID             string
	TradeLocked      bool
	Addons           []string
	BuyerCountryCode string
}
