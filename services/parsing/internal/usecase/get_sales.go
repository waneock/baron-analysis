package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"skinbaron-analyzer/pkg/messaging/jobs"
	"skinbaron-analyzer/services/parsing/internal/client/baron"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"time"
)

var (
	ErrClientGetSales    = fmt.Errorf("skinbaron api call error")
	ErrBodyWithoutOffers = fmt.Errorf("skinbaron api call response: body without offers")
	ErrInsertIntoDb      = fmt.Errorf("error when trying to insert into database")
)

type OffersClient interface {
	GetSales(ctx context.Context, afterSaleID string) (*baron.GetSalesResponse, error)
}

type OffersRepo interface {
	CreateMany(ctx context.Context, offers []domain.Offer) error
}

type JobsRepo interface {
	UpdateStatus(ctx context.Context, id string, status jobs.SyncJobStatus) error
}

type SyncOffers struct {
	client      OffersClient
	offersRepo  OffersRepo
	jobsRepo    JobsRepo
	logger      *slog.Logger
	afterSaleID string
}

func NewSyncOffers(client OffersClient, offersRepo OffersRepo, jobsRepo JobsRepo, logger *slog.Logger) *SyncOffers {
	return &SyncOffers{
		client:      client,
		offersRepo:  offersRepo,
		jobsRepo:    jobsRepo,
		logger:      logger,
		afterSaleID: "0",
	}
}

func (uc *SyncOffers) Execute(ctx context.Context, jobID string) {
	uc.jobsRepo.UpdateStatus(ctx, jobID, jobs.SyncJobStatusRunning)

	if err := uc.doSync(ctx); err != nil {
		uc.logger.Error("sync offers do sync",
			"error", err)
		uc.jobsRepo.UpdateStatus(ctx, jobID, jobs.SyncJobStatusFailed)
		return
	}

	uc.jobsRepo.UpdateStatus(ctx, jobID, jobs.SyncJobStatusDone)
}

func (uc *SyncOffers) doSync(ctx context.Context) error {
	cnt := 0
	for {
		cnt += 1
		uc.logger.Info("sync offers: ",
			"iteration", cnt)

		salesResponse, err := uc.client.GetSales(ctx, uc.afterSaleID)
		if err != nil {
			uc.logger.Error("sync offers",
				"error", err)
			return ErrClientGetSales
		}

		if salesResponse == nil {
			uc.logger.Info("sync offers: sales response nil")
			return nil
		}

		offers := clientResponseToOffer(*salesResponse)

		if offers == nil || len(*offers) == 0 {
			uc.logger.Warn("sync offers: empty offers")
			return ErrBodyWithoutOffers
		}

		err = uc.offersRepo.CreateMany(ctx, *offers)
		if err != nil {
			uc.logger.Error("sync offers",
				"error", err)
			return ErrInsertIntoDb
		}

		uc.afterSaleID = (*offers)[len(*offers)-1].ID
	}
}

func clientResponseToOffer(salesResponse baron.GetSalesResponse) *[]domain.Offer {
	var offers []domain.Offer
	for _, res := range salesResponse.Response {
		offer := domain.Offer{
			ID:               res.ID,
			Price:            int(res.Price * 100),
			Commission:       int(res.Commission * 100),
			Tax:              res.Tax,
			ClassID:          res.ClassID,
			InstanceID:       res.InstanceID,
			AppID:            res.AppID,
			ContextID:        res.ContextID,
			AssetID:          res.AssetID,
			Name:             res.Name,
			OfferID:          res.OfferID,
			State:            res.State,
			EscrowEndDate:    time.Unix(res.EscrowEndDate, 0),
			ListTime:         time.Unix(res.ListTime, 0),
			LastUpdated:      time.Unix(res.LastUpdated, 0),
			Wear:             int(res.Wear * 100),
			TxID:             res.TxID,
			TradeLocked:      res.TradeLocked,
			Addons:           res.Addons,
			BuyerCountryCode: res.BuyerCountryCode,
		}
		offers = append(offers, offer)
	}
	return &offers
}
