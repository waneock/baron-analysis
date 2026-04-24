package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"skinbaron-analyzer/services/parsing/internal/repository"
	"strings"
)

type ItemWearSaleRepo struct {
	db *sql.DB
}

func NewItemWearSaleRepo(db *sql.DB) *ItemWearSaleRepo {
	return &ItemWearSaleRepo{
		db: db,
	}
}

func (iw *ItemWearSaleRepo) CreateMany(ctx context.Context, items []domain.ItemWearSale) error {
	if len(items) == 0 {
		return nil
	}

	var (
		stringBuilder strings.Builder
		args          []any
	)

	stringBuilder.WriteString(`
		INSERT INTO item_wear_sales (
			item_wear_id,
			price,
			wear_value,
			sold_on
		) VALUES
	`)

	argPos := 1
	for i, item := range items {
		if i > 0 {
			stringBuilder.WriteString(",")
		}

		stringBuilder.WriteString(fmt.Sprintf(
			"($%d,$%d,$%d,$%d)",
			argPos, argPos+1, argPos+2, argPos+3,
		))

		args = append(args,
			item.WearID,
			item.Price,
			item.Wear,
			item.DateSold,
		)

		argPos += 4
	}

	stringBuilder.WriteString(`
		ON CONFLICT (item_wear_id, sold_on, price, wear_value) DO NOTHING;
	`)

	ctx, cancel := context.WithTimeout(ctx, repository.QueryRequestTimeout)
	defer cancel()

	_, err := iw.db.ExecContext(ctx, stringBuilder.String(), args...)
	return err
}
