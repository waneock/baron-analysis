package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"strings"
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

	_, err := o.db.ExecContext(ctx, stringBuilder.String(), args...)
	return err
}
