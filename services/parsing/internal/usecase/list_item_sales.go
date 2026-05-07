package usecase

import (
	"context"
	"skinbaron-analyzer/services/parsing/internal/domain"
)

type ItemSalesRepo interface {
	ListSales(ctx context.Context, filter domain.ListItemSalesFilter) ([]domain.ItemSales, error)
	ListSalesStats(ctx context.Context, filter domain.ListItemSalesStatsFilter) ([]domain.ItemSalesStats, error)
}

type ItemSalesService struct {
	itemSalesRepo ItemSalesRepo
}

func NewItemSalesService(itemSalesRepo ItemSalesRepo) *ItemSalesService {
	return &ItemSalesService{
		itemSalesRepo: itemSalesRepo,
	}
}

func (s *ItemSalesService) ListSales(ctx context.Context, filter domain.ListItemSalesFilter) (*domain.ListItemSalesOutput, error) {
	if filter.Offset < 0 {
		filter.Offset = listOffsetFielDefault
	}

	if filter.Limit < 0 || filter.Limit > listLimitFieldMax {
		filter.Limit = listLimitFieldMax
	}

	var output domain.ListItemSalesOutput
	items, err := s.itemSalesRepo.ListSales(ctx, filter)
	if err != nil {
		return nil, err
	}

	output.Items = items
	output.Limit = filter.Limit
	output.Offset = filter.Offset

	return &output, nil
}

func (s *ItemSalesService) ListSalesStat(ctx context.Context, filter domain.ListItemSalesStatsFilter) (*domain.ListItemSalesStatsOutput, error) {
	if filter.Offset < 0 {
		filter.Offset = listOffsetFielDefault
	}

	if filter.Limit < 0 || filter.Limit > listLimitFieldMax {
		filter.Limit = listLimitFieldMax
	}

	var output domain.ListItemSalesStatsOutput
	items, err := s.itemSalesRepo.ListSalesStats(ctx, filter)
	if err != nil {
		return nil, err
	}

	output.Items = items
	output.Limit = filter.Limit
	output.Offset = filter.Offset

	return &output, nil
}
