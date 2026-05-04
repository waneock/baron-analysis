package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"skinbaron-analyzer/pkg/messaging/jobs"
	"skinbaron-analyzer/services/parsing/internal/domain"
)

const (
	limitDefaultValue = 500
)

var (
	ErrCountItems     = fmt.Errorf("error when trying to get items count")
	ErrListItems      = fmt.Errorf("error when trying to get items")
	ErrGetNewestSales = fmt.Errorf("error when trying to get newest sales")
	ErrCreateItems    = fmt.Errorf("error when trying to write items into db")
)

type BaronClient interface {
	GetNewestSales(ctx context.Context, itemName string) (*[]domain.GetNewestSalesOut, error)
}

type MarketSyncSourceRepo interface {
	Count(ctx context.Context) (int, error)
	List(ctx context.Context, limit, offset int) (*[]domain.ItemWearSource, error)
}

type ItemWearSalesRepo interface {
	CreateMany(ctx context.Context, items []domain.ItemWearSale) error
}

type SyncItemPrices struct {
	itemRepo             ItemsRepo
	marketSyncSourceRepo MarketSyncSourceRepo
	itemWearSalesRepo    ItemWearSalesRepo
	jobsRepo             JobsRepo
	baronClient          BaronClient
	log                  *slog.Logger
}

func NewSyncItemPrices(itemRepo ItemsRepo,
	marketSyncSourceRepo MarketSyncSourceRepo,
	itemWearSalesRepo ItemWearSalesRepo,
	jobsRepo JobsRepo,
	baronClient BaronClient,
	log *slog.Logger) *SyncItemPrices {
	return &SyncItemPrices{
		itemRepo:             itemRepo,
		marketSyncSourceRepo: marketSyncSourceRepo,
		itemWearSalesRepo:    itemWearSalesRepo,
		jobsRepo:             jobsRepo,
		baronClient:          baronClient,
		log:                  log,
	}
}

func (uc *SyncItemPrices) Execute(ctx context.Context, jobID string) {
	uc.jobsRepo.UpdateStatus(ctx, jobID, jobs.SyncJobStatusRunning)

	if err := uc.doSync(ctx); err != nil {
		uc.log.Error("sync items do sync",
			"error", err)
		uc.jobsRepo.UpdateStatus(ctx, jobID, jobs.SyncJobStatusFailed)
	}

	uc.jobsRepo.UpdateStatus(ctx, jobID, jobs.SyncJobStatusDone)
}

func (uc *SyncItemPrices) doSync(ctx context.Context) error {
	total, err := uc.marketSyncSourceRepo.Count(ctx)
	if err != nil {
		uc.log.Error("error when trying to count the items",
			"error", err)
		return ErrCountItems
	}

	for i := 3214; i < total; i += limitDefaultValue {
		uc.log.Info("sync item prices",
			"index", i,
			"total", total)

		limit := i + limitDefaultValue
		if limit > total {
			limit = total - i
		}

		items, err := uc.marketSyncSourceRepo.List(ctx, limit, i)
		if err != nil {
			uc.log.Error("market sync source repo list",
				"err", err)

			return ErrListItems
		}

		for j, item := range *items {
			itemName := fmt.Sprintf("%s (%s)", item.Name, item.Wear)
			newestSales, err := uc.baronClient.GetNewestSales(ctx, itemName)
			if err != nil {
				uc.log.Error("get newest sales",
					"err", err)
				return ErrGetNewestSales
			}

			itemWearSales := newestSalesOutToItemWearSale(newestSales, item.WearID)

			fmt.Println("item index: ", i+j)
			// fmt.Println("items:", itemWearSales)

			err = uc.itemWearSalesRepo.CreateMany(ctx, itemWearSales)
			if err != nil {
				uc.log.Error("error when trying to add data into item wears table",
					"err", err)
				return ErrCreateItems
			}
		}
	}

	return nil
}

func newestSalesOutToItemWearSale(input *[]domain.GetNewestSalesOut, wearID int) []domain.ItemWearSale {
	items := make([]domain.ItemWearSale, 0, len(*input))
	for _, item := range *input {
		newItem := domain.ItemWearSale{
			WearID:   wearID,
			Price:    item.Price,
			Wear:     item.Wear,
			DateSold: item.DateSold,
		}

		items = append(items, newItem)
	}

	return items
}
