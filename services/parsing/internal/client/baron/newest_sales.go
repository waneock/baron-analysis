package baron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"time"
)

const (
	getNewestSales30DaysEndpoint = "GetNewestSales30Days"
	statTrakDefaultValue         = "false"
	souvenirDefaultValue         = "false"
)

type GetNewestPayload struct {
	Apikey   string `json:"apikey"`
	ItemName string `json:"itemName"`
	StatTrak string `json:"statTrak"`
	Souvenir string `json:"souvenir"`
}

type GetNewestEntry struct {
	ItemName string  `json:"itemName"`
	Price    float64 `json:"price"`
	Wear     float64 `json:"wear"`
	DateSold string  `json:"dateSold"`
}

type GetNewestResponse struct {
	NewestSales30Days []GetNewestEntry `json:"newestSales30Days"`
}

func (b *BaronClient) GetNewestSales(ctx context.Context, itemName string) (*[]domain.GetNewestSalesOut, error) {
	getNewestSalesPayload := createGetNewestSalesBody(itemName, b.apiKey)

	bodyBytes, err := json.Marshal(getNewestSalesPayload)
	if err != nil {
		return nil, fmt.Errorf("baron: offer: marshal request body %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req, err := b.createRequest(ctx, http.MethodPost, getNewestSales30DaysEndpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	var resp GetNewestResponse
	if err := b.do(req, &resp); err != nil {
		return nil, err
	}

	items := getNewestSalesResponseToOut(resp)

	return &items, nil
}

func createGetNewestSalesBody(itemName, apiKey string) GetNewestPayload {
	return GetNewestPayload{
		Apikey:   apiKey,
		ItemName: itemName,
		StatTrak: statTrakDefaultValue,
		Souvenir: souvenirDefaultValue,
	}
}

func getNewestSalesResponseToOut(input GetNewestResponse) []domain.GetNewestSalesOut {
	items := make([]domain.GetNewestSalesOut, 0, len(input.NewestSales30Days))
	for _, item := range input.NewestSales30Days {
		dateSold, _ := time.Parse("2006-01-02", item.DateSold) // TODO: don't forget to handle the error
		newItem := domain.GetNewestSalesOut{
			ItemName: item.ItemName,
			Price:    item.Price,
			Wear:     item.Wear,
			DateSold: dateSold,
		}
		items = append(items, newItem)
	}
	return items
}
