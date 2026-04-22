package parsinggrpc

import (
	pb "skinbaron-analyzer/proto/parsing/v1"
	"skinbaron-analyzer/services/reporting/internal/domain"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	appIDDefaultValue int64 = 730
)

func listOffersInputToPbRequest(input domain.ListOffersInput) *pb.ListOffersRequest {
	var state *int64
	if input.State != nil {
		state = input.State
	}

	var nameQuery *string
	if input.NameQuery != nil {
		nameQuery = input.NameQuery
	}

	var minPrice *float64
	if input.MinPrice != nil {
		minPrice = input.MinPrice
	}

	var maxPrice *float64
	if input.MaxPrice != nil {
		maxPrice = input.MaxPrice
	}

	var listTime *timestamppb.Timestamp
	if input.ListTime != nil {
		listTime = toProtoTime(input.ListTime)
	}

	var lastUpdated *timestamppb.Timestamp
	if input.LastUpdated != nil {
		lastUpdated = toProtoTime(input.LastUpdated)
	}

	var sortBy *string
	if input.SortBy != nil {
		sortBy = input.SortBy
	}

	var sortOrder *string
	if input.SortOrder != nil {
		sortOrder = input.SortOrder
	}

	return &pb.ListOffersRequest{
		Limit:       input.Limit,
		Offset:      input.Offset,
		State:       state,
		AppId:       &appIDDefaultValue,
		NameQuery:   nameQuery,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		ListTime:    listTime,
		LastUpdated: lastUpdated,
		SortBy:      sortBy,
		SortOrder:   sortOrder,
	}
}

func pbResponseToListOffersOut(input *pb.ListOffersResponse) *domain.ListOffersOutput {
	items := make([]domain.Offer, 0, len(input.Items))
	for _, item := range input.Items {
		newOffer := domain.Offer{
			ID:         item.GetId(),
			Price:      item.GetPrice(),
			Commission: item.GetCommission(),
			Name:       item.GetName(),
			State:      strconv.FormatInt(item.GetState(), 10),
			Wear:       item.GetWear(),
			Addons:     item.GetAddons(),
		}

		listTime := toTime(item.GetListTime())
		if listTime != nil {
			newOffer.ListTime = *listTime
		}

		lastUpdated := toTime(item.GetLastUpdated())
		if lastUpdated != nil {
			newOffer.LastUpdated = *lastUpdated
		}

		items = append(items, newOffer)
	}
	return &domain.ListOffersOutput{
		Items:  items,
		Total:  input.GetTotal(),
		Limit:  input.GetLimit(),
		Offset: input.GetOffset(),
	}
}

func toTime(t *timestamppb.Timestamp) *time.Time {
	if t == nil {
		return nil
	}

	val := (*t).AsTime()
	return &val
}

func toProtoTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(*t)
}
