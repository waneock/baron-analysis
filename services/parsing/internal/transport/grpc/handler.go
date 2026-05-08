package grpc

import (
	"context"
	pb "skinbaron-analyzer/proto/parsing/v1"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"skinbaron-analyzer/services/parsing/internal/usecase"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListOffersUseCase interface {
	Execute(ctx context.Context, input usecase.ListOffersInput) (*usecase.ListOfferResult, error)
}

type ListItemSalesUseCase interface {
	ListSales(ctx context.Context, filter domain.ListItemSalesFilter) (*domain.ListItemSalesOutput, error)
	ListSalesStat(ctx context.Context, filter domain.ListItemSalesStatsFilter) (*domain.ListItemSalesStatsOutput, error)
}

type Handler struct {
	pb.UnimplementedParsingServiceServer

	listOffersUC    ListOffersUseCase
	listItemSalesUC ListItemSalesUseCase
}

func NewHandler(listOffersUC ListOffersUseCase, listItemSalesUC ListItemSalesUseCase) *Handler {
	return &Handler{
		listOffersUC:    listOffersUC,
		listItemSalesUC: listItemSalesUC,
	}
}

func (h *Handler) ListOffers(ctx context.Context, req *pb.ListOffersRequest) (*pb.ListOffersResponse, error) {
	input := usecase.ListOffersInput{
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	}

	if req.State != nil {
		val := int(req.GetState())
		input.State = &val
	}

	if req.AppId != nil {
		val := int(req.GetAppId())
		input.AppID = &val
	}

	if req.NameQuery != nil {
		val := req.GetNameQuery()
		input.NameQuery = &val
	}

	if req.MinPrice != nil {
		val := int(req.GetMinPrice() * 100)
		input.MinPrice = &val
	}

	if req.MaxPrice != nil {
		val := int(req.GetMaxPrice() * 100)
		input.MaxPrice = &val
	}

	if req.ListTime != nil {
		val := toTime(req.GetListTime())
		input.ListTime = val
	}

	if req.LastUpdated != nil {
		val := toTime(req.GetLastUpdated())
		input.LastUpdated = val
	}

	if req.SortBy != nil {
		val := req.GetSortBy()
		input.SortBy = &val
	}

	if req.SortOrder != nil {
		val := req.GetSortOrder()
		input.SortOrder = &val
	}

	result, err := h.listOffersUC.Execute(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list offers: %v", err)
	}

	items := make([]*pb.Offer, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toProtoOffer(item))
	}

	return &pb.ListOffersResponse{
		Items:  items,
		Total:  result.Total,
		Limit:  int64(result.Limit),
		Offset: int64(result.Offset),
	}, nil

}

func toProtoOffer(offer domain.Offer) *pb.Offer {
	return &pb.Offer{
		Id:               offer.ID,
		Price:            float64(offer.Price) / 100.00,
		Commission:       float64(offer.Commission) / 100.00,
		Tax:              float64(offer.Tax),
		ClassId:          offer.ClassID,
		InstanceId:       offer.InstanceID,
		AppId:            int64(offer.AppID),
		ContextId:        offer.ContextID,
		AssetId:          offer.AssetID,
		Name:             offer.Name,
		OfferId:          offer.OfferID,
		State:            int64(offer.State),
		EscrowEndDate:    toProtoTime(&offer.EscrowEndDate),
		LastUpdated:      toProtoTime(&offer.LastUpdated),
		ListTime:         toProtoTime(&offer.ListTime),
		Wear:             float64(offer.Wear) / 100.00,
		Txid:             offer.TxID,
		TradeLocked:      offer.TradeLocked,
		Addons:           offer.Addons,
		BuyerCountryCode: offer.BuyerCountryCode,
	}
}

func toProtoTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(*t)
}

func toTime(t *timestamppb.Timestamp) *time.Time {
	if t == nil {
		return nil
	}

	val := (*t).AsTime()
	return &val
}

func (h *Handler) ListItemSales(ctx context.Context, req *pb.ListItemSalesRequest) (*pb.ListItemSalesResponse, error) {
	filter := domain.ListItemSalesFilter{
		Offset: req.GetOffset(),
		Limit:  req.GetLimit(),
	}

	if req.ItemNameQuery != nil && req.GetItemNameQuery() != "" {
		val := req.GetItemNameQuery()
		filter.ItemNameQuery = &val
	}

	if req.WearName != nil && req.GetWearName() != "" {
		val := req.GetWearName()
		filter.WearName = &val
	}

	if req.ItemWearId != nil {
		val := req.GetItemWearId()
		filter.ItemWearID = &val
	}

	if req.MinPrice != nil {
		val := req.GetMinPrice()
		filter.MinPrice = &val
	}

	if req.MaxPrice != nil {
		val := req.GetMaxPrice()
		filter.MaxPrice = &val
	}

	if req.SoldFrom != nil && req.GetSoldFrom() != "" {
		val := req.GetSoldFrom()
		filter.SoldFrom = stringToTime(val)
	}

	if req.SoldTo != nil && req.GetSoldTo() != "" {
		val := req.GetSoldTo()
		filter.SoldTo = stringToTime(val)
	}

	ucOutput, err := h.listItemSalesUC.ListSales(ctx, filter)
	if err != nil {
		return nil, err
	}

	return itemSalesOutputToResponse(*ucOutput), nil
}

func itemSalesOutputToResponse(input domain.ListItemSalesOutput) *pb.ListItemSalesResponse {
	pbItems := make([]*pb.ItemSale, 0, len(input.Items))
	for _, item := range input.Items {
		pbItem := pb.ItemSale{
			ItemName:  item.ItemName,
			WearName:  item.WearName,
			Price:     item.Price,
			WearValue: item.WearValue,
			SoldOn:    timeToString(item.SoldOn),
		}
		pbItems = append(pbItems, &pbItem)
	}

	return &pb.ListItemSalesResponse{
		Items:  pbItems,
		Limit:  input.Limit,
		Offset: input.Offset,
	}
}

func stringToTime(stringTime string) *time.Time {
	val, err := time.Parse("2006-01-02", stringTime)
	if err != nil {
		return nil
	}
	return &val
}

func timeToString(val time.Time) string {
	return val.Format("2006-01-02")
}

func (h *Handler) ListItemSalesStats(ctx context.Context, req *pb.ListItemSalesStatsRequest) (*pb.ListItemSalesStatsResponse, error) {
	filter := domain.ListItemSalesStatsFilter{
		Limit:  req.GetLimit(),
		Offset: req.GetOffset(),
	}

	if req.ItemNameQuery != nil && req.GetItemNameQuery() != "" {
		val := req.GetItemNameQuery()
		filter.ItemNameQuery = &val
	}

	if req.WearName != nil && req.GetWearName() != "" {
		val := req.GetWearName()
		filter.WearName = &val
	}

	if req.MinPrice != nil {
		val := req.GetMinPrice()
		filter.MinPrice = &val
	}

	if req.MaxPrice != nil {
		val := req.GetMaxPrice()
		filter.MaxPrice = &val
	}

	if req.SoldFrom != nil && req.GetSoldFrom() != "" {
		val := req.GetSoldFrom()
		filter.SoldFrom = stringToTime(val)
	}

	if req.SoldTo != nil && req.GetSoldTo() != "" {
		val := req.GetSoldTo()
		filter.SoldTo = stringToTime(val)
	}

	if req.MinSalesCount != nil {
		val := req.GetMinSalesCount()
		filter.MinSalesCount = &val
	}

	ucOutput, err := h.listItemSalesUC.ListSalesStat(ctx, filter)
	if err != nil {
		return nil, err
	}

	return itemSalesStateOutputToResponse(*ucOutput), nil
}

func itemSalesStateOutputToResponse(input domain.ListItemSalesStatsOutput) *pb.ListItemSalesStatsResponse {
	pbItems := make([]*pb.ItemSalesStats, 0, len(input.Items))
	for _, item := range input.Items {
		pbItem := &pb.ItemSalesStats{
			ItemId:      item.ItemID,
			ItemName:    item.ItemName,
			ItemWearId:  item.ItemWearID,
			WearName:    item.WearName,
			SalesCount:  item.SalesCount,
			AvgPrice:    item.AvgPrice,
			MedianPrice: item.MedianPrice,
			MinPrice:    item.MinPrice,
			MaxPrice:    item.MaxPrice,
			SoldPrices:  item.SoldPrices,
			FirstSoldOn: timeToString(item.FirstSoldOn),
			LastSoldOn:  timeToString(item.LastSoldOn),
		}
		pbItems = append(pbItems, pbItem)
	}
	return &pb.ListItemSalesStatsResponse{
		Items:  pbItems,
		Limit:  input.Limit,
		Offset: input.Offset,
	}
}
