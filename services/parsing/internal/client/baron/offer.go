package baron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	getSalesEndpoint = "GetSales"
	cs2AppID         = 730
	itemsPerPage     = 50
	sortOrderAsc     = 0
	requestType      = 0
)

type GetSalesPayload struct {
	ApiKey       string `json:"apikey"`
	Type         int    `json:"type"`
	AppID        int    `json:"appid"`
	AfterSaleID  string `json:"after_saleid"`
	ItemsPerPage int    `json:"items_per_page"`
	SortOrder    int    `json:"sort_order"`
}

type GetSalesOffer struct {
	ID               string   `json:"id"`
	Price            float64  `json:"price"`
	Commission       float64  `json:"commission"`
	Tax              int      `json:"tax"`
	ClassID          string   `json:"classid"`
	InstanceID       string   `json:"instanceid"`
	AppID            int      `json:"appid"`
	ContextID        string   `json:"contextid"`
	AssetID          string   `json:"assetid"`
	Name             string   `json:"name"`
	OfferID          string   `json:"offerid"`
	State            int      `json:"state"`
	EscrowEndDate    int64    `json:"escrow_end_date"`
	ListTime         int64    `json:"list_time"`
	LastUpdated      int64    `json:"last_updated"`
	Wear             float64  `json:"wear"`
	TxID             string   `json:"txid"`
	TradeLocked      bool     `json:"trade_locked"`
	Addons           []string `json:"addons"`
	BuyerCountryCode string   `json:"buyer_country_code"`
}

type GetSalesResponse struct {
	Response []GetSalesOffer `json:"response"`
}

func (b *BaronClient) GetSales(ctx context.Context, afterSaleID string) (*GetSalesResponse, error) {
	getSalesPayload := createGetSalesOffer(b.apiKey, afterSaleID)

	fmt.Println("get sales payload: ", getSalesPayload)

	bodyBytes, err := json.Marshal(getSalesPayload)
	if err != nil {
		return nil, fmt.Errorf("baron: offer: marshal request body %w", err)
	}

	fmt.Println("get sales payload: ", string(bodyBytes))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := b.createRequest(ctx, http.MethodPost, getSalesEndpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	var resp GetSalesResponse
	if err := b.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, err
}

func createGetSalesOffer(apiKey, afterSaleID string) *GetSalesPayload {
	return &GetSalesPayload{
		ApiKey:       apiKey,
		Type:         requestType,
		AppID:        cs2AppID,
		AfterSaleID:  afterSaleID,
		ItemsPerPage: itemsPerPage,
		SortOrder:    sortOrderAsc,
	}
}
