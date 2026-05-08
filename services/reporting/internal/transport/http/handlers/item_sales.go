package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"skinbaron-analyzer/services/reporting/internal/domain"
	"skinbaron-analyzer/services/reporting/internal/transport/http/render"
	"strconv"
)

const (
	paramWearName      = "wear_name"
	paramWearID        = "wear_id"
	paramSoldFrom      = "sold_from"
	paramSoldTo        = "sold_to"
	paramMinSalesCount = "min_sales_count"
)

type ItemSalesService interface {
	ListItemSales(ctx context.Context, input domain.ListItemSalesInput) (*domain.ListItemSalesOutput, error)
	ListItemSalesStats(ctx context.Context, input domain.ListItemSalesStatInput) (*domain.ListItemSalesStatOutput, error)
}

type ItemSalesHandler struct {
	itemSalesService ItemSalesService
	logger           *slog.Logger
}

func NewItemSalesHandler(itemSalesService ItemSalesService, logger *slog.Logger) *ItemSalesHandler {
	return &ItemSalesHandler{
		itemSalesService: itemSalesService,
		logger:           logger,
	}
}

type ItemSale struct {
	ItemName  string  `json:"item_name"`
	WearName  string  `json:"wear_name"`
	Price     float64 `json:"price"`
	WearValue float64 `json:"wear_value"`
	SoldOn    string  `json:"sold_on"`
}

type ListItemSalesOutput struct {
	Items  []ItemSale `json:"items"`
	Limit  int64      `json:"limit"`
	Offset int64      `json:"offset"`
}

func (h *ItemSalesHandler) ListItemSales(w http.ResponseWriter, r *http.Request) {
	payload, err := parseListItemSalesQuery(r)
	if err != nil {
		render.BadRequestErr(w)
		return
	}

	ctx := r.Context()
	listItemSales, err := h.itemSalesService.ListItemSales(ctx, *payload)
	if err != nil {
		render.InternalServerErr(w)
		return
	}

	items := svcOutToItemSalesOut(*listItemSales)
	if err := render.OK(w, items); err != nil {
		render.InternalServerErr(w)
		return
	}
}

func parseListItemSalesQuery(r *http.Request) (*domain.ListItemSalesInput, error) {
	q := r.URL.Query()

	limit, err := strconv.ParseInt(q.Get(paramLimit), 10, 64)
	if err != nil {
		return nil, err
	}

	offset, err := strconv.ParseInt(q.Get(paramOffset), 10, 64)
	if err != nil {
		return nil, err
	}

	var itemNameQuery *string
	if raw := q.Get(paramNameQuery); raw != "" {
		itemNameQuery = &raw
	}

	var wearName *string
	if raw := q.Get(paramWearName); raw != "" {
		wearName = &raw
	}

	var itemWearID *int64
	if raw := q.Get(paramWearID); raw != "" {
		val, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, err
		}
		itemWearID = &val
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

	var soldFrom *string
	if raw := q.Get(paramSoldFrom); raw != "" {
		soldFrom = &raw
	}

	var soldTo *string
	if raw := q.Get(paramSoldTo); raw != "" {
		soldTo = &raw
	}

	return &domain.ListItemSalesInput{
		Limit:         limit,
		Offset:        offset,
		ItemNameQuery: itemNameQuery,
		WearName:      wearName,
		ItemWearID:    itemWearID,
		MinPrice:      minPrice,
		MaxPrice:      maxPrice,
		SoldFrom:      soldFrom,
		SoldTo:        soldTo,
	}, nil
}

func svcOutToItemSalesOut(input domain.ListItemSalesOutput) ListItemSalesOutput {
	items := make([]ItemSale, 0, len(input.Items))
	for _, item := range input.Items {
		newItem := ItemSale{
			ItemName:  item.ItemName,
			WearName:  item.WearName,
			Price:     item.Price,
			WearValue: item.WearValue,
			SoldOn:    item.SoldOn,
		}
		items = append(items, newItem)
	}
	return ListItemSalesOutput{
		Items:  items,
		Limit:  input.Limit,
		Offset: input.Offset,
	}
}

type ItemSalesStats struct {
	ItemID      string    `json:"item_id"`
	ItemName    string    `json:"item_name"`
	ItemWearID  int64     `json:"item_wear_id"`
	WearName    string    `json:"wear_name"`
	SalesCount  int64     `json:"sales_count"`
	AvgPrice    float64   `json:"avg_price"`
	MedianPrice float64   `json:"median_price"`
	MinPrice    float64   `json:"min_price"`
	MaxPrice    float64   `json:"max_price"`
	SoldPrices  []float64 `json:"sold_prices"`
	FirstSoldOn string    `json:"first_sold_on"`
	LastSoldOn  string    `json:"last_sold_on"`
}

type ListItemSalesStatOutput struct {
	Items  []ItemSalesStats `json:"items"`
	Limit  int64            `json:"limit"`
	Offset int64            `json:"offset"`
}

func (h *ItemSalesHandler) ListItemSalesStats(w http.ResponseWriter, r *http.Request) {
	payload, err := parseListItemSalesStatsQuery(r)
	if err != nil {
		render.BadRequestErr(w)
		return
	}

	ctx := r.Context()
	listItemSalesStat, err := h.itemSalesService.ListItemSalesStats(ctx, *payload)
	if err != nil {
		render.InternalServerErr(w)
		return
	}

	items := svcOutToItemSalesStatOut(*listItemSalesStat)
	if err := render.OK(w, items); err != nil {
		render.InternalServerErr(w)
		return
	}
}

func parseListItemSalesStatsQuery(r *http.Request) (*domain.ListItemSalesStatInput, error) {
	q := r.URL.Query()

	limit, err := strconv.ParseInt(q.Get(paramLimit), 10, 64)
	if err != nil {
		return nil, err
	}

	offset, err := strconv.ParseInt(q.Get(paramOffset), 10, 64)
	if err != nil {
		return nil, err
	}

	var itemNameQuery *string
	if raw := q.Get(paramNameQuery); raw != "" {
		itemNameQuery = &raw
	}

	var wearName *string
	if raw := q.Get(paramWearName); raw != "" {
		wearName = &raw
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

	var soldFrom *string
	if raw := q.Get(paramSoldFrom); raw != "" {
		soldFrom = &raw
	}

	var soldTo *string
	if raw := q.Get(paramSoldTo); raw != "" {
		soldTo = &raw
	}

	var minSalesCount *int64
	temp, err := strconv.ParseInt(q.Get(paramMinSalesCount), 10, 64)
	if err == nil {
		minSalesCount = &temp
	}

	return &domain.ListItemSalesStatInput{
		Limit:         limit,
		Offset:        offset,
		ItemNameQuery: itemNameQuery,
		WearName:      wearName,
		MinPrice:      minPrice,
		MaxPrice:      maxPrice,
		SoldFrom:      soldFrom,
		SoldTo:        soldTo,
		MinSalesCount: minSalesCount,
	}, nil
}

func svcOutToItemSalesStatOut(input domain.ListItemSalesStatOutput) ListItemSalesStatOutput {
	items := make([]ItemSalesStats, 0, len(input.Items))
	for _, item := range input.Items {
		newItem := ItemSalesStats{
			ItemID:      item.ItemID,
			ItemName:    item.ItemName,
			ItemWearID:  item.ItemWearID,
			WearName:    item.WearName,
			SalesCount:  item.SalesCount,
			AvgPrice:    item.AvgPrice,
			MedianPrice: item.MedianPrice,
			MinPrice:    item.MinPrice,
			MaxPrice:    item.MaxPrice,
			SoldPrices:  item.SoldPrices,
			FirstSoldOn: item.FirstSoldOn,
			LastSoldOn:  item.LastSoldOn,
		}
		items = append(items, newItem)
	}

	return ListItemSalesStatOutput{
		Items:  items,
		Limit:  input.Limit,
		Offset: input.Offset,
	}
}
