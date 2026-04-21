package http

import (
	"context"
	"net/http"
	pb "skinbaron-analyzer/proto/parsing/v1"
	"skinbaron-analyzer/services/reporting/internal/transport/http/render"
	ucsvc "skinbaron-analyzer/services/reporting/internal/usecase"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	paramLimit       = "limit"
	paramOffset      = "offset"
	paramAppID       = "app_id"
	paramState       = "state"
	paramNameQuery   = "name_query"
	paramMinPrice    = "min_price"
	paramMaxPrice    = "max_price"
	paramListTime    = "list_time"
	paramLastUpdated = "last_updated"
	paramSortBy      = "sort_by"
	paramSortOrder   = "sort_order"
)

type SyncOffersService interface {
	Execute(ctx context.Context) (*pb.SyncOffersResponse, error)
}

type ListOffersService interface {
	Execute(ctx context.Context, input ucsvc.ListOffersInput) (*pb.ListOffersResponse, error)
}

type OffersHandler struct {
	syncOffers SyncOffersService
	listOffers ListOffersService
}

func NewOffersHandler(syncOffers SyncOffersService, listOffers ListOffersService) *OffersHandler {
	return &OffersHandler{
		syncOffers: syncOffers,
		listOffers: listOffers,
	}
}

func (h *OffersHandler) SyncOffers(w http.ResponseWriter, r *http.Request) {

}

type ListOffersPayload struct {
	Limit       int        `json:"limit"`
	Offset      int        `json:"offset"`
	AppID       *int       `json:"app_id"`
	State       *int       `json:"state"`
	NameQuery   *string    `json:"name_query"`
	MinPrice    *float64   `json:"min_price"`
	MaxPrice    *float64   `json:"max_price"`
	ListTime    *time.Time `json:"list_time"`
	LastUpdated *time.Time `json:"last_updated"`
	SortBy      *string    `json:"sort_by"`
	SortOrder   *string    `json:"sort_order"`
}

type ListOffer struct {
	ID               string    `json:"id"`
	Price            float64   `json:"price"`
	Commission       float64   `json:"commission"`
	Tax              float64   `json:"tax"`
	ClassID          string    `json:"class_id"`
	InstanceID       string    `json:"instance_id"`
	AppID            int64     `json:"app_id"`
	ContextID        string    `json:"context_id"`
	AssetID          string    `json:"asset_id"`
	Name             string    `json:"name"`
	OfferID          string    `json:"offer_id"`
	State            int64     `json:"state"`
	EscrowEndDate    time.Time `json:"escrow_end_date"`
	ListTime         time.Time `json:"list_time"`
	LastUpdated      time.Time `json:"last_updated"`
	Wear             float64   `json:"wear"`
	TxID             string    `json:"txid"`
	TradeLocked      bool      `json:"trade_locked"`
	Addons           []string  `json:"addons"`
	BuyerCountryCode string    `json:"buyer_country_code"`
}

type ListOfferOutput struct {
	Items []ListOffer `json:"items"`
}

func (h *OffersHandler) ListOffers(w http.ResponseWriter, r *http.Request) {
	payload, err := parseListOffersQuery(r)
	if err != nil {
		render.BadRequestErr(w)
		return
	}

	ctx := r.Context()
	svcInput := listOffersPayloadToSvcInput(*payload)

	offersResponse, err := h.listOffers.Execute(ctx, *svcInput)
	if err != nil {
		render.InternalServerErr(w)
		return
	}

	items := pbListOffersResToListOffersOut(offersResponse)
	if err := render.OK(w, items); err != nil {
		render.InternalServerErr(w)
		return
	}
}

func listOffersPayloadToSvcInput(input ListOffersPayload) *ucsvc.ListOffersInput {
	return &ucsvc.ListOffersInput{
		Offset:      input.Offset,
		Limit:       input.Limit,
		AppID:       input.AppID,
		State:       input.State,
		NameQuery:   input.NameQuery,
		MinPrice:    input.MinPrice,
		MaxPrice:    input.MaxPrice,
		ListTime:    input.ListTime,
		LastUpdated: input.LastUpdated,
		SortBy:      input.SortBy,
		SortOrder:   input.SortOrder,
	}
}

func pbListOffersResToListOffersOut(input *pb.ListOffersResponse) *ListOfferOutput {
	items := make([]ListOffer, 0, len(input.Items))
	for _, item := range input.Items {
		newOffer := ListOffer{
			ID:               item.GetId(),
			Price:            item.GetPrice(),
			Commission:       item.GetCommission(),
			Tax:              item.GetTax(),
			ClassID:          item.GetClassId(),
			InstanceID:       item.GetInstanceId(),
			AppID:            item.GetAppId(),
			ContextID:        item.GetContextId(),
			AssetID:          item.GetAssetId(),
			Name:             item.GetName(),
			OfferID:          item.GetOfferId(),
			State:            item.GetState(),
			Wear:             item.GetWear(),
			TxID:             item.GetTxid(),
			TradeLocked:      item.GetTradeLocked(),
			Addons:           item.GetAddons(),
			BuyerCountryCode: item.GetBuyerCountryCode(),
		}

		escrowEndDate := toTime(item.GetEscrowEndDate())
		if escrowEndDate != nil {
			newOffer.EscrowEndDate = *escrowEndDate
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
	return &ListOfferOutput{
		Items: items,
	}
}

func toTime(t *timestamppb.Timestamp) *time.Time {
	if t == nil {
		return nil
	}

	val := (*t).AsTime()
	return &val
}

func parseListOffersQuery(r *http.Request) (*ListOffersPayload, error) {
	q := r.URL.Query()

	limit, err := strconv.Atoi(q.Get(paramLimit))
	if err != nil {
		return nil, err
	}

	offset, err := strconv.Atoi(q.Get(paramOffset))
	if err != nil {
		return nil, err
	}

	var appID *int
	if raw := q.Get(paramAppID); raw != "" {
		val, err := strconv.Atoi(raw)
		if err != nil {
			return nil, err
		}
		appID = &val
	}

	var state *int
	if raw := q.Get(paramState); raw != "" {
		val, err := strconv.Atoi(raw)
		if err != nil {
			return nil, err
		}
		state = &val
	}

	var nameQuery *string
	if raw := q.Get(paramNameQuery); raw != "" {
		nameQuery = &raw
	}

	var minPrice *float64
	if raw := q.Get(paramMinPrice); raw != "" {
		val, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, err
		}
		minPrice = &val
	}

	var maxPrice *float64
	if raw := q.Get(paramMaxPrice); raw != "" {
		val, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, err
		}
		maxPrice = &val
	}

	// todo add the rest of parameters

	return &ListOffersPayload{
		Limit:       limit,
		Offset:      offset,
		AppID:       appID,
		State:       state,
		NameQuery:   nameQuery,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		ListTime:    nil,
		LastUpdated: nil,
		SortBy:      nil,
		SortOrder:   nil,
	}, nil

}
