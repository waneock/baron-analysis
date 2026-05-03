package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"skinbaron-analyzer/pkg/messaging/jobs"
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
	jobsRepo      JobsRepo
	log           *slog.Logger
	jsonPath      string
}

func NewItemsService(itemsRepo ItemsRepo,
	itemWearsRepo ItemWearsRepo,
	jobsRepo JobsRepo,
	log *slog.Logger,
	jsonPath string) *ItemsService {
	return &ItemsService{
		itemsRepo:     itemsRepo,
		itemWearsRepo: itemWearsRepo,
		jobsRepo:      jobsRepo,
		log:           log,
		jsonPath:      jsonPath,
	}
}

func (uc *ItemsService) Execute(ctx context.Context, jobID string) {
	uc.jobsRepo.UpdateStatus(ctx, jobID, jobs.SyncJobStatusRunning)

	if err := uc.doSync(ctx); err != nil {
		uc.log.Error("sync items do sync",
			"error", err)
		uc.jobsRepo.UpdateStatus(ctx, jobID, jobs.SyncJobStatusFailed)
	}

	uc.jobsRepo.UpdateStatus(ctx, jobID, jobs.SyncJobStatusDone)
}

func (uc *ItemsService) doSync(ctx context.Context) error {
	items, err := localjson.ReadItemJSON(uc.jsonPath)
	if err != nil {
		return fmt.Errorf("error when trying to read item json: %w", err)
	}

	for i := 0; i < len(*items); i += 100 {
		uc.log.Info("processing items",
			"index", i,
			"from", len(*items))
		end := i + 100
		if end > len(*items) {
			end = len(*items)
		}

		batch := (*items)[i:end]

		itemRows := uc.itemsToItemRows(batch)
		if err := uc.itemsRepo.CreateMany(ctx, itemRows); err != nil {
			return fmt.Errorf("error when tryign to write into items table: %w", err)
		}

		itemWearRows := uc.itemsToItemWearRows(batch)
		if err := uc.itemWearsRepo.CreateMany(ctx, itemWearRows); err != nil {
			return fmt.Errorf("error when trying to write into item_wears table: %w", err)
		}
	}

	return nil
}

func (uc *ItemsService) itemsToItemRows(items []domain.Item) []domain.ItemRow {
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

func (uc *ItemsService) itemsToItemWearRows(items []domain.Item) []domain.ItemWearRow {
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
