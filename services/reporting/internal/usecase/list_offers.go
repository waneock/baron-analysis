package usecase

import (
	"context"
	"skinbaron-analyzer/services/reporting/internal/domain"
)

type ListOffersClient interface {
	ListOffers(ctx context.Context, input domain.ListOffersInput) (*domain.ListOffersOutput, error)
}

type ListOffers struct {
	client ListOffersClient
}

func NewListOffers(client ListOffersClient) *ListOffers {
	return &ListOffers{
		client: client,
	}
}

func (uc *ListOffers) Execute(ctx context.Context, input domain.ListOffersInput) (*domain.ListOffersOutput, error) {
	offers, err := uc.client.ListOffers(ctx, input)
	if err != nil {
		return nil, err
	}
	for i := range offers.Items {
		offers.Items[i].State = stringifyOfferState(offers.Items[i].State)
	}
	return offers, nil
}

func stringifyOfferState(state string) string {
	switch state {
	case "2":
		return "Available"
	case "4":
		return "Sold"
	case "7":
		return "Canceled"
	default:
		return state
	}
}
