package usecase

import (
	"context"
	pb "skinbaron-analyzer/proto/parsing/v1"
)

type SyncOffersClient interface {
	SyncOffers(ctx context.Context) (*pb.SyncOffersResponse, error)
}

type SyncOffers struct {
	client SyncOffersClient
}

func NewSyncOffers(client SyncOffersClient) *SyncOffers {
	return &SyncOffers{
		client: client,
	}
}

func (uc *SyncOffers) Execute(ctx context.Context) (*pb.SyncOffersResponse, error) {
	return uc.client.SyncOffers(ctx)
}
