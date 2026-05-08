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
	var nameQuery *string
	if input.NameQuery != nil && *input.NameQuery != "" {
		nameQuery = input.NameQuery
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
	if input.SortBy != nil && *input.SortBy != "" {
		sortBy = input.SortBy
	}

	var sortOrder *string
	if input.SortOrder != nil && *input.SortOrder != "" {
		sortOrder = input.SortOrder
	}

	return &pb.ListOffersRequest{
		Limit:       input.Limit,
		Offset:      input.Offset,
		State:       input.State,
		AppId:       &appIDDefaultValue,
		NameQuery:   nameQuery,
		MinPrice:    input.MaxPrice,
		MaxPrice:    input.MinPrice,
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

func listItemSalesInputToPbRequest(input domain.ListItemSalesInput) pb.ListItemSalesRequest {
	var itemNameQuery *string
	if input.ItemNameQuery != nil && *input.ItemNameQuery != "" {
		itemNameQuery = input.ItemNameQuery
	}

	var wearName *string
	if input.WearName != nil && *input.WearName != "" {
		wearName = input.WearName
	}

	var soldFrom *string
	if input.SoldFrom != nil && *input.SoldFrom != "" {
		soldFrom = input.SoldFrom
	}

	var soldTo *string
	if input.SoldTo != nil && *input.SoldTo != "" {
		soldTo = input.SoldTo
	}

	return pb.ListItemSalesRequest{
		Limit:         input.Limit,
		Offset:        input.Offset,
		ItemNameQuery: itemNameQuery,
		WearName:      wearName,
		ItemWearId:    input.ItemWearID,
		MinPrice:      input.MinPrice,
		MaxPrice:      input.MaxPrice,
		SoldFrom:      soldFrom,
		SoldTo:        soldTo,
	}
}

func pbResponseToListItemSalesOut(input *pb.ListItemSalesResponse) domain.ListItemSalesOutput {
	itemSales := make([]domain.ItemSale, 0, len(input.Items))
	for _, item := range input.Items {
		newItem := domain.ItemSale{
			ItemName:  item.GetItemName(),
			WearName:  item.GetWearName(),
			Price:     item.GetPrice(),
			WearValue: item.GetWearValue(),
			SoldOn:    item.GetSoldOn(),
		}
		itemSales = append(itemSales, newItem)
	}

	return domain.ListItemSalesOutput{
		Items:  itemSales,
		Limit:  input.Limit,
		Offset: input.Offset,
	}
}

func listItemSalesStatsInputToPbRequest(input domain.ListItemSalesStatInput) pb.ListItemSalesStatsRequest {
	var itemNameQuery *string
	if input.ItemNameQuery != nil && *input.ItemNameQuery != "" {
		itemNameQuery = input.ItemNameQuery
	}

	var wearName *string
	if input.WearName != nil && *input.WearName != "" {
		wearName = input.WearName
	}

	var soldFrom *string
	if input.SoldFrom != nil && *input.SoldFrom != "" {
		soldFrom = input.SoldFrom
	}

	var soldTo *string
	if input.SoldTo != nil && *input.SoldTo != "" {
		soldTo = input.SoldTo
	}

	return pb.ListItemSalesStatsRequest{
		Limit:         input.Limit,
		Offset:        input.Offset,
		ItemNameQuery: itemNameQuery,
		WearName:      wearName,
		MinPrice:      input.MinPrice,
		MaxPrice:      input.MaxPrice,
		SoldFrom:      soldFrom,
		SoldTo:        soldTo,
		MinSalesCount: input.MinSalesCount,
	}
}

func pbResponseToListItemSalesStatsOut(input *pb.ListItemSalesStatsResponse) domain.ListItemSalesStatOutput {
	items := make([]domain.ItemSalesStats, 0, len(input.Items))

	for _, item := range input.Items {
		newItem := domain.ItemSalesStats{
			ItemID:      item.GetItemId(),
			ItemName:    item.GetItemName(),
			ItemWearID:  item.GetItemWearId(),
			WearName:    item.GetWearName(),
			SalesCount:  item.GetSalesCount(),
			AvgPrice:    item.GetAvgPrice(),
			MedianPrice: item.GetMedianPrice(),
			MinPrice:    item.GetMinPrice(),
			MaxPrice:    item.GetMaxPrice(),
			SoldPrices:  item.GetSoldPrices(),
			FirstSoldOn: item.GetFirstSoldOn(),
			LastSoldOn:  item.GetLastSoldOn(),
		}
		items = append(items, newItem)
	}

	return domain.ListItemSalesStatOutput{
		Items:  items,
		Limit:  input.Limit,
		Offset: input.Offset,
	}
}
