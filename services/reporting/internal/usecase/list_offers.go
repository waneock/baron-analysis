package usecase

import (
	"context"
	pb "skinbaron-analyzer/proto/parsing/v1"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListOffersClient interface {
	ListOffers(ctx context.Context, req *pb.ListOffersRequest) (*pb.ListOffersResponse, error)
}

type ListOffers struct {
	client ListOffersClient
}

func NewListOffers(client ListOffersClient) *ListOffers {
	return &ListOffers{
		client: client,
	}
}

type ListOffersInput struct {
	Limit  int
	Offset int

	AppID       *int
	State       *int
	NameQuery   *string
	MinPrice    *float64
	MaxPrice    *float64
	ListTime    *time.Time
	LastUpdated *time.Time

	SortBy    *string
	SortOrder *string
}

func (uc *ListOffers) Execute(ctx context.Context, input ListOffersInput) (*pb.ListOffersResponse, error) {
	req := &pb.ListOffersRequest{
		Limit:  int64(input.Limit),
		Offset: int64(input.Offset),
	}

	if input.AppID != nil {
		val := int64(*input.AppID)
		req.AppId = &val
	}

	if input.State != nil {
		val := int64(*input.State)
		req.State = &val
	}

	if input.NameQuery != nil {
		req.NameQuery = input.NameQuery
	}

	if input.MinPrice != nil {
		req.MinPrice = input.MinPrice
	}

	if input.MaxPrice != nil {
		req.MaxPrice = input.MaxPrice
	}

	if input.ListTime != nil {
		req.ListTime = timestamppb.New(*input.ListTime)
	}

	if input.LastUpdated != nil {
		req.LastUpdated = timestamppb.New(*input.LastUpdated)
	}

	if input.SortBy != nil {
		req.SortBy = input.SortBy
	}

	if input.SortOrder != nil {
		req.SortOrder = input.SortOrder
	}

	return uc.client.ListOffers(ctx, req)

}
