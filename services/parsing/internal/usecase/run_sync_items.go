package usecase

import (
	"context"
	"log/slog"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"skinbaron-analyzer/services/parsing/internal/source/localjson"
)

type ItemsRepo interface {
	CreateMany(ctx context.Context, items []domain.ItemRow) error
}

type ItemWearsRepo interface {
	CreateMany(ctx context.Context, wears []domain.ItemWearRow) error
}

type ItemsService struct {
	itemsRepo     ItemsRepo
	itemWearsRepo ItemWearsRepo
	log           *slog.Logger
	jsonPath      string
}

func NewItemsService(itemsRepo ItemsRepo, itemWearsRepo ItemWearsRepo, log *slog.Logger, jsonPath string) *ItemsService {
	return &ItemsService{
		itemsRepo:     itemsRepo,
		itemWearsRepo: itemWearsRepo,
		log:           log,
		jsonPath:      jsonPath,
	}
}

func (s *ItemsService) Execute(ctx context.Context) {
	items, err := localjson.ReadItemJSON(s.jsonPath)
	if err != nil {
		s.log.Error("error when trying to read item json",
			"path", s.jsonPath,
			"error", err)
	}

	for i := 0; i < len(*items); i += 100 {
		s.log.Info("processing items",
			"index", i,
			"from", len(*items))
		end := i + 100
		if end > len(*items) {
			end = len(*items)
		}

		batch := (*items)[i:end]

		itemRows := s.itemsToItemRows(batch)
		if err := s.itemsRepo.CreateMany(ctx, itemRows); err != nil {
			s.log.Error("error when trying to write into items table",
				"error", err)
		}

		itemWearRows := s.itemsToItemWearRows(batch)
		if err := s.itemWearsRepo.CreateMany(ctx, itemWearRows); err != nil {
			s.log.Error("error when trying to write into item_wears table",
				"error", err)
		}
	}
}

func (s *ItemsService) itemsToItemRows(items []domain.Item) []domain.ItemRow {
	itemRows := make([]domain.ItemRow, 0, len(items))

	for _, item := range items {
		row := domain.ItemRow{
			ID:   item.ID,
			Name: item.Name,
		}

		itemRows = append(itemRows, row)
	}

	return itemRows
}

func (s *ItemsService) itemsToItemWearRows(items []domain.Item) []domain.ItemWearRow {
	itemWearRows := make([]domain.ItemWearRow, 0, len(items))

	for _, item := range items {

		for _, wear := range item.Wears {
			row := domain.ItemWearRow{
				ID:   item.ID,
				Name: wear,
			}

			itemWearRows = append(itemWearRows, row)
		}
	}

	return itemWearRows
}
