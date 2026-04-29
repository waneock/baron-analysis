package http

import (
	"context"
	"net/http"
	"skinbaron-analyzer/services/reporting/internal/domain"
	"skinbaron-analyzer/services/reporting/internal/transport/http/render"
	"strconv"
	"time"
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
	Execute(ctx context.Context, jobType domain.SyncJobType) (string, error)
}

type ListOffersService interface {
	Execute(ctx context.Context, input domain.ListOffersInput) (*domain.ListOffersOutput, error)
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

type SyncOffersOutput struct {
	JobID string `json:"job_id"`
}

func (h *OffersHandler) SyncOffers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobID, err := h.syncOffers.Execute(ctx, domain.SyncJobTypeSyncOffers)
	if err != nil {
		render.InternalServerErr(w)
		return
	}

	var syncOffersOutput SyncOffersOutput
	syncOffersOutput.JobID = jobID
	if err := render.OK(w, syncOffersOutput); err != nil {
		render.InternalServerErr(w)
		return
	}
}

type ListOffersPayload struct {
	Limit       int64      `json:"limit"`
	Offset      int64      `json:"offset"`
	State       *int64     `json:"state"`
	NameQuery   *string    `json:"name_query"`
	MinPrice    *float64   `json:"min_price"`
	MaxPrice    *float64   `json:"max_price"`
	ListTime    *time.Time `json:"list_time"`
	LastUpdated *time.Time `json:"last_updated"`
	SortBy      *string    `json:"sort_by"`
	SortOrder   *string    `json:"sort_order"`
}

type ListOffer struct {
	ID          string    `json:"id"`
	Price       float64   `json:"price"`
	Commission  float64   `json:"commission"`
	Name        string    `json:"name"`
	State       string    `json:"state"`
	ListTime    time.Time `json:"list_time"`
	LastUpdated time.Time `json:"last_updated"`
	Wear        float64   `json:"wear"`
	Addons      []string  `json:"addons"`
}

type ListOfferOutput struct {
	Items  []ListOffer `json:"items"`
	Total  int64       `json:"total"`
	Limit  int64       `json:"limit"`
	Offset int64       `json:"offset"`
}

func (h *OffersHandler) ListOffers(w http.ResponseWriter, r *http.Request) {
	payload, err := parseListOffersQuery(r)
	if err != nil {
		render.BadRequestErr(w)
		return
	}

	ctx := r.Context()
	svcInput := listOffersPayloadToSvcInput(*payload)

	offersResponse, err := h.listOffers.Execute(ctx, svcInput)
	if err != nil {
		render.InternalServerErr(w)
		return
	}

	items := svcOutToListOffersOut(*offersResponse)
	if err := render.OK(w, items); err != nil {
		render.InternalServerErr(w)
		return
	}
}

func parseListOffersQuery(r *http.Request) (*ListOffersPayload, error) {
	q := r.URL.Query()

	limit, err := strconv.ParseInt(q.Get(paramLimit), 10, 64)
	if err != nil {
		return nil, err
	}

	offset, err := strconv.ParseInt(q.Get(paramOffset), 10, 64)
	if err != nil {
		return nil, err
	}

	var state *int64
	if raw := q.Get(paramState); raw != "" {
		val, err := strconv.ParseInt(raw, 10, 64)
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

func listOffersPayloadToSvcInput(input ListOffersPayload) domain.ListOffersInput {
	return domain.ListOffersInput{
		Limit:       input.Limit,
		Offset:      input.Offset,
		AppID:       nil, // this value is constant and set inside client
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

func svcOutToListOffersOut(input domain.ListOffersOutput) *ListOfferOutput {
	items := make([]ListOffer, 0, len(input.Items))

	for _, item := range input.Items {
		newOffer := ListOffer{
			ID:          item.ID,
			Price:       item.Price,
			Commission:  item.Commission,
			Name:        item.Name,
			State:       item.State,
			ListTime:    item.ListTime,
			LastUpdated: item.LastUpdated,
			Wear:        item.Wear,
			Addons:      item.Addons,
		}

		items = append(items, newOffer)
	}

	return &ListOfferOutput{
		Items:  items,
		Total:  input.Total,
		Limit:  input.Limit,
		Offset: input.Offset,
	}
}
