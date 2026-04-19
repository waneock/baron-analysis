package usecase

import (
	"context"
	"log/slog"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"skinbaron-analyzer/services/parsing/internal/repository"
	"time"
)

const (
	listLimitFieldMax     = 100
	listOffsetFielDefault = 0
)

type ListRepo interface {
	List(ctx context.Context, filter repository.OfferFilter) ([]domain.Offer, error)
	Count(ctx context.Context, filter repository.OfferFilter) (int64, error)
}

type ListOffersService struct {
	repo   ListRepo
	logger *slog.Logger
}

type ListOfferResult struct {
	Items  []domain.Offer
	Total  int64
	Limit  int
	Offset int
}

type ListOffersInput struct {
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

func NewListOfferService(listRepo ListRepo, logger *slog.Logger) *ListOffersService {
	return &ListOffersService{
		repo:   listRepo,
		logger: logger,
	}
}

func (uc *ListOffersService) Execute(ctx context.Context, input ListOffersInput) (*ListOfferResult, error) {
	if input.Limit < 0 {
		input.Limit = 0
	}

	if input.Limit > listLimitFieldMax {
		input.Limit = listLimitFieldMax
	}

	if input.Offset < 0 {
		input.Offset = listOffsetFielDefault
	}

	filter := listOffersInputToOffersFilter(input)

	items, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	total, err := uc.repo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &ListOfferResult{
		Items:  items,
		Total:  total,
		Offset: filter.Offset,
		Limit:  filter.Limit,
	}, nil

}

func listOffersInputToOffersFilter(input ListOffersInput) repository.OfferFilter {
	return repository.OfferFilter{
		Limit:       input.Limit,
		Offset:      input.Offset,
		AppID:       input.AppID,
		State:       input.State,
		NameQuery:   input.NameQuery,
		MinPrice:    input.MinPrice,
		MaxPrice:    input.MaxPrice,
		ListTime:    input.ListTime,
		LastUpdated: input.LastUpdated,
		SortBy:      input.SortBy,
		SortOrder:   input.SortOrder,
	}
}
