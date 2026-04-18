package usecase

import (
	"context"
	"log/slog"
	"skinbaron-analyzer/services/parsing/internal/client/baron"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"time"
)

type OffersClient interface {
	GetSales(ctx context.Context, afterSaleID string) (*baron.GetSalesResponse, error)
}

type OffersRepo interface {
	CreateMany(ctx context.Context, offers []domain.Offer) error
}

type GetSalesService struct {
	client      OffersClient
	repo        OffersRepo
	logger      *slog.Logger
	afterSaleID string
}

func NewGetSalesService(client OffersClient, repo OffersRepo, logger *slog.Logger) *GetSalesService {
	return &GetSalesService{
		client:      client,
		repo:        repo,
		logger:      logger,
		afterSaleID: "0",
	}
}

func (uc *GetSalesService) SyncOffers(ctx context.Context) {
	cnt := 0
	for {
		cnt += 1
		uc.logger.Info("sync offers: ",
			"iteration", cnt)

		salesResponse, err := uc.client.GetSales(ctx, uc.afterSaleID)
		if err != nil {
			uc.logger.Error("sync offers",
				"error", err)
		}

		if salesResponse == nil {
			uc.logger.Warn("sync offers: sales response nil")
			return
		}

		offers := clientResponseToOffer(*salesResponse)

		if offers == nil || len(*offers) == 0 {
			uc.logger.Warn("sync offers: empty offers")
			return
		}

		err = uc.repo.CreateMany(ctx, *offers)
		if err != nil {
			uc.logger.Error("sync offers",
				"error", err)
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
