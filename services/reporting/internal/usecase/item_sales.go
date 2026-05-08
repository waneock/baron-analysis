package usecase

import (
	"context"
	"skinbaron-analyzer/services/reporting/internal/domain"
)

type ItemSalesClient interface {
	ListItemSales(ctx context.Context, input domain.ListItemSalesInput) (*domain.ListItemSalesOutput, error)
	ListItemSalesStats(ctx context.Context, input domain.ListItemSalesStatInput) (*domain.ListItemSalesStatOutput, error)
}

type ItemSales struct {
	itemSalesClient ItemSalesClient
}

func NewItemSales(itemSalesClient ItemSalesClient) *ItemSales {
	return &ItemSales{
		itemSalesClient: itemSalesClient,
	}
}

func (uc *ItemSales) ListItemSales(ctx context.Context, input domain.ListItemSalesInput) (*domain.ListItemSalesOutput, error) {
	return uc.itemSalesClient.ListItemSales(ctx, input)
}

func (uc *ItemSales) ListItemSalesStats(ctx context.Context, input domain.ListItemSalesStatInput) (*domain.ListItemSalesStatOutput, error) {
	return uc.itemSalesClient.ListItemSalesStats(ctx, input)
}
