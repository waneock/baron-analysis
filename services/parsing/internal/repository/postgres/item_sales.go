package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"skinbaron-analyzer/services/parsing/internal/repository"
	"strings"
)

type ItemSalesRepo struct {
	db *sql.DB
}

func NewItemSalesRepo(db *sql.DB) *ItemSalesRepo {
	return &ItemSalesRepo{
		db: db,
	}
}

func (i *ItemSalesRepo) ListSales(ctx context.Context, filter domain.ListItemSalesFilter) ([]domain.ItemSales, error) {
	var (
		stringBuilder strings.Builder
		args          []any
	)

	stringBuilder.WriteString(`
		SELECT
			i.name AS item_name,
			iw.name AS wear_name,
			iws.price,
			iws.wear_value,
			iws.sold_on
		FROM item_wear_sales iws
		JOIN item_wears iw ON iw.id = iws.item_wear_id
		JOIN items i ON i.id = iw.item_id
		WHERE 1 = 1
	`)

	argPos := 1

	if filter.ItemNameQuery != nil && *filter.ItemNameQuery != "" {
		stringBuilder.WriteString(fmt.Sprintf(" AND i.name ILIKE $%d", argPos))
		args = append(args, "%"+*filter.ItemNameQuery+"%")
		argPos++
	}

	if filter.WearName != nil && *filter.WearName != "" {
		stringBuilder.WriteString(fmt.Sprintf(" AND iw.name = $%d", argPos))
		args = append(args, *filter.WearName)
		argPos++
	}

	if filter.ItemWearID != nil {
		stringBuilder.WriteString(fmt.Sprintf(" AND iw.id = $%d", argPos))
		args = append(args, *filter.ItemWearID)
		argPos++
	}

	if filter.MinPrice != nil {
		stringBuilder.WriteString(fmt.Sprintf(" AND iws.price >= $%d", argPos))
		args = append(args, *filter.MinPrice)
		argPos++
	}

	if filter.MaxPrice != nil {
		stringBuilder.WriteString(fmt.Sprintf(" AND iws.price <= $%d", argPos))
		args = append(args, *filter.MaxPrice)
		argPos++
	}

	if filter.SoldFrom != nil {
		stringBuilder.WriteString(fmt.Sprintf(" AND iws.sold_on >= $%d", argPos))
		args = append(args, *filter.SoldFrom)
		argPos++
	}

	if filter.SoldTo != nil {
		stringBuilder.WriteString(fmt.Sprintf(" AND iws.sold_on <= $%d", argPos))
		args = append(args, *filter.SoldTo)
		argPos++
	}

	stringBuilder.WriteString(fmt.Sprintf(`
		ORDER BY iws.sold_on DESC, i.name, iw.name
		LIMIT $%d OFFSET $%d
	`, argPos, argPos+1))

	args = append(args, filter.Limit, filter.Offset)

	ctx, cancel := context.WithTimeout(ctx, repository.QueryRequestTimeout)
	defer cancel()

	rows, err := i.db.QueryContext(ctx, stringBuilder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.ItemSales, 0)

	for rows.Next() {
		var item domain.ItemSales

		if err := rows.Scan(
			&item.ItemName,
			&item.WearName,
			&item.Price,
			&item.WearValue,
			&item.SoldOn,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (i *ItemSalesRepo) ListSalesStats(ctx context.Context, filter domain.ListItemSalesStatsFilter) ([]domain.ItemSalesStats, error) {
	var (
		stringBuilder strings.Builder
		args          []any
	)

	stringBuilder.WriteString(`
		SELECT
			i.id AS item_id,
			i.name AS item_name,
			iw.id AS item_wear_id,
			iw.name AS wear_name,

			COUNT(*) AS sales_count,
			AVG(iws.price)::float8 AS avg_price,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY iws.price)::float8 AS median_price,
			MIN(iws.price)::float8 AS min_price,
			MAX(iws.price)::float8 AS max_price,

			COALESCE(
				json_agg(iws.price::float8 ORDER BY iws.sold_on DESC, iws.id DESC),
				'[]'::json
			) AS sold_prices,

			MIN(iws.sold_on) AS first_sold_on,
			MAX(iws.sold_on) AS last_sold_on
		FROM item_wear_sales iws
		JOIN item_wears iw ON iw.id = iws.item_wear_id
		JOIN items i ON i.id = iw.item_id
		WHERE 1 = 1
	`)

	argPos := 1

	if filter.ItemNameQuery != nil && *filter.ItemNameQuery != "" {
		stringBuilder.WriteString(fmt.Sprintf(" AND i.name ILIKE $%d", argPos))
		args = append(args, "%"+*filter.ItemNameQuery+"%")
		argPos++
	}

	if filter.WearName != nil && *filter.WearName != "" {
		stringBuilder.WriteString(fmt.Sprintf(" AND iw.name = $%d", argPos))
		args = append(args, *filter.WearName)
		argPos++
	}

	if filter.MinPrice != nil {
		stringBuilder.WriteString(fmt.Sprintf(" AND iws.price >= $%d", argPos))
		args = append(args, *filter.MinPrice)
		argPos++
	}

	if filter.MaxPrice != nil {
		stringBuilder.WriteString(fmt.Sprintf(" AND iws.price <= $%d", argPos))
		args = append(args, *filter.MaxPrice)
		argPos++
	}

	if filter.SoldFrom != nil {
		stringBuilder.WriteString(fmt.Sprintf(" AND iws.sold_on >= $%d", argPos))
		args = append(args, *filter.SoldFrom)
		argPos++
	}

	if filter.SoldTo != nil {
		stringBuilder.WriteString(fmt.Sprintf(" AND iws.sold_on <= $%d", argPos))
		args = append(args, *filter.SoldTo)
		argPos++
	}

	stringBuilder.WriteString(`
		GROUP BY
			i.id,
			i.name,
			iw.id,
			iw.name
	`)

	if filter.MinSalesCount != nil {
		stringBuilder.WriteString(fmt.Sprintf(" HAVING COUNT(*) >= $%d", argPos))
		args = append(args, *filter.MinSalesCount)
		argPos++
	}

	stringBuilder.WriteString(fmt.Sprintf(`
		ORDER BY
			sales_count DESC,
			avg_price DESC
		LIMIT $%d OFFSET $%d
	`, argPos, argPos+1))

	args = append(args, filter.Limit, filter.Offset)

	ctx, cancel := context.WithTimeout(ctx, repository.QueryRequestTimeout)
	defer cancel()

	rows, err := i.db.QueryContext(ctx, stringBuilder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.ItemSalesStats, 0)

	for rows.Next() {
		var (
			item          domain.ItemSalesStats
			soldPricesRaw []byte
		)

		if err := rows.Scan(
			&item.ItemID,
			&item.ItemName,
			&item.ItemWearID,
			&item.WearName,
			&item.SalesCount,
			&item.AvgPrice,
			&item.MedianPrice,
			&item.MinPrice,
			&item.MaxPrice,
			&soldPricesRaw,
			&item.FirstSoldOn,
			&item.LastSoldOn,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(soldPricesRaw, &item.SoldPrices); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
