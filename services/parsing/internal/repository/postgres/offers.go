package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"skinbaron-analyzer/services/parsing/internal/repository"
	"strings"

	"github.com/lib/pq"
)

type OffersRepo struct {
	db *sql.DB
}

func NewOffersRepo(db *sql.DB) *OffersRepo {
	return &OffersRepo{
		db: db,
	}
}

func (o *OffersRepo) CreateMany(ctx context.Context, offers []domain.Offer) error {
	if len(offers) == 0 {
		return nil
	}

	var (
		stringBuilder strings.Builder
		args          []any
	)

	stringBuilder.WriteString(`
		INSERT INTO orders (
			id, price, commission, tax, classid, instanceid, appid, contextid,
			assetid, name, offerid, state, escrow_end_date, list_time,
			last_updated, wear, txid, trade_locked, addons, buyer_country_code
		) VALUES
	`)

	argPos := 1
	for i, offer := range offers {
		if i > 0 {
			stringBuilder.WriteString(",")
		}

		stringBuilder.WriteString(fmt.Sprintf(
			"($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			argPos, argPos+1, argPos+2, argPos+3, argPos+4,
			argPos+5, argPos+6, argPos+7, argPos+8, argPos+9,
			argPos+10, argPos+11, argPos+12, argPos+13, argPos+14,
			argPos+15, argPos+16, argPos+17, argPos+18, argPos+19,
		))

		args = append(args,
			offer.ID,
			offer.Price,
			offer.Commission,
			offer.Tax,
			offer.ClassID,
			offer.InstanceID,
			offer.AppID,
			offer.ContextID,
			offer.AssetID,
			offer.Name,
			offer.OfferID,
			offer.State,
			offer.EscrowEndDate,
			offer.ListTime,
			offer.LastUpdated,
			offer.Wear,
			offer.TxID,
			offer.TradeLocked,
			offer.Addons,
			offer.BuyerCountryCode,
		)

		argPos += 20
	}

	stringBuilder.WriteString(`
		ON CONFLICT (id) DO UPDATE SET
				price = EXCLUDED.price,
				commission = EXCLUDED.commission,
				tax = EXCLUDED.tax,
				state = EXCLUDED.state,
				last_updated = EXCLUDED.last_updated
	`)

	ctx, cancel := context.WithTimeout(ctx, repository.QueryRequestTimeout)
	defer cancel()

	_, err := o.db.ExecContext(ctx, stringBuilder.String(), args...)
	return err
}

func (o *OffersRepo) List(ctx context.Context, filter repository.OfferFilter) ([]domain.Offer, error) {
	var args []any

	query := `
		SELECT
				id,
				price,
				commission,
				tax,
				classid,
				instanceid,
				appid,
				contextid,
				assetid,
				name,
				offerid,
				state,
				escrow_end_date,
				list_time,
				last_updated,
				wear,
				txid,
				trade_locked,
				addons,
				buyer_country_code
		FROM orders
			WHERE 1=1
		`
	argPos := 1

	addFiltersToQuery(filter, &query, &args, &argPos)

	sortBy := normalizeSortBy(filter.SortBy)
	sortOrder := normalizeSortOrder(filter.SortOrder)

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filter.Limit, filter.Offset)

	ctx, cancel := context.WithTimeout(ctx, repository.QueryRequestTimeout)
	defer cancel()

	rows, err := o.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []domain.Offer

	for rows.Next() {
		var offer domain.Offer

		err := rows.Scan(
			&offer.ID,
			&offer.Price,
			&offer.Commission,
			&offer.Tax,
			&offer.ClassID,
			&offer.InstanceID,
			&offer.AppID,
			&offer.ContextID,
			&offer.AssetID,
			&offer.Name,
			&offer.OfferID,
			&offer.State,
			&offer.EscrowEndDate,
			&offer.ListTime,
			&offer.LastUpdated,
			&offer.Wear,
			&offer.TxID,
			&offer.TradeLocked,
			pq.Array(&offer.Addons),
			&offer.BuyerCountryCode,
		)

		if err != nil {
			return nil, err
		}

		offers = append(offers, offer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return offers, nil
}

func (o *OffersRepo) Count(ctx context.Context, filter repository.OfferFilter) (int64, error) {
	query := `
		SELECT
			COUNT(*)
		FROM
			orders
		WHERE 1=1	
	`

	var args []any
	argPos := 1
	addFiltersToQuery(filter, &query, &args, &argPos)

	ctx, cancel := context.WithTimeout(ctx, repository.QueryRequestTimeout)
	defer cancel()

	var total int64
	err := o.db.QueryRowContext(ctx, query, args...).Scan(&total)
	return total, err
}

func normalizeSortBy(sortBy *string) string {
	if sortBy == nil {
		return "list_time"
	}

	switch *sortBy {
	case "price":
		return "price"
	case "list_time":
		return "list_time"
	case "last_updated":
		return "last_updated"
	default:
		return "list_time"
	}
}

func normalizeSortOrder(sortOrder *string) string {
	if sortOrder == nil {
		return "DESC"
	}

	switch *sortOrder {
	case "asc":
		return "ASC"
	default:
		return "DESC"
	}
}

func addFiltersToQuery(filter repository.OfferFilter, query *string, args *[]any, argPos *int) {
	if filter.AppID != nil {
		(*query) += fmt.Sprintf(" AND appid = $%d", *argPos)
		*args = append(*args, *filter.AppID)
		(*argPos)++
	}

	if filter.State != nil {
		(*query) += fmt.Sprintf(" AND state = $%d", *argPos)
		*args = append(*args, *filter.State)
		(*argPos)++
	}

	if filter.NameQuery != nil && *filter.NameQuery != "" {
		(*query) += fmt.Sprintf(" AND name ILIKE $%d", *argPos)
		*args = append(*args, "%"+*filter.NameQuery+"%")
		(*argPos)++
	}

	if filter.MinPrice != nil {
		(*query) += fmt.Sprintf(" AND price >= $%d", *argPos)
		*args = append(*args, *filter.MinPrice)
		(*argPos)++
	}

	if filter.MaxPrice != nil {
		(*query) += fmt.Sprintf(" AND price <= $%d", *argPos)
		*args = append(*args, *filter.MaxPrice)
		(*argPos)++
	}

	if filter.ListTime != nil {
		(*query) += fmt.Sprintf(" AND list_time >= $%d", *argPos)
		*args = append(*args, *filter.ListTime)
		(*argPos)++
	}

	if filter.LastUpdated != nil {
		(*query) += fmt.Sprintf(" AND last_updated >= $%d", *argPos)
		*args = append(*args, *filter.LastUpdated)
		(*argPos)++
	}
}
