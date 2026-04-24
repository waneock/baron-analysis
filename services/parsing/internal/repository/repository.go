package repository

import (
	"context"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"time"
)

const (
	QueryRequestTimeout = 5 * time.Second
)

type OfferFilter struct {
	Limit  int
	Offset int

	AppID       *int
	State       *int
	NameQuery   *string
	MinPrice    *int
	MaxPrice    *int
	ListTime    *time.Time
	LastUpdated *time.Time

	SortBy    *string
	SortOrder *string
}

type OffersRepository interface {
	CreateMany(ctx context.Context, offers []domain.Offer) error
	List(ctx context.Context, filter OfferFilter) ([]domain.Offer, error)
	Count(ctx context.Context, filter OfferFilter) (int64, error)
}

type ItemsRepository interface {
	CreateMany(ctx context.Context, items []domain.ItemRow) error
}

type ItemWearsRepository interface {
	CreateMany(ctx context.Context, wears []domain.ItemWearRow) error
}

type MarketSyncSourceRepository interface {
	Count(ctx context.Context) (int, error)
	List(ctx context.Context, limit, offset int) (*[]domain.ItemWearSource, error)
}

type ItemWearSaleRepository interface {
	CreateMany(ctx context.Context, items []domain.ItemWearSale) error
}

type Repo struct {
	OffersRepository           OffersRepository
	ItemsRepository            ItemsRepository
	ItemWearsRepository        ItemWearsRepository
	MarketSyncSourceRepository MarketSyncSourceRepository
	ItemWearSaleRepository     ItemWearSaleRepository
}
