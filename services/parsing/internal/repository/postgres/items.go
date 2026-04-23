package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"skinbaron-analyzer/services/parsing/internal/repository"
	"strings"
)

type ItemsRepo struct {
	db *sql.DB
}

func NewItemsRepo(db *sql.DB) *ItemsRepo {
	return &ItemsRepo{
		db: db,
	}
}

func (i *ItemsRepo) CreateMany(ctx context.Context, items []domain.ItemRow) error {
	if len(items) == 0 {
		return nil
	}

	var (
		stringBuilder strings.Builder
		args          []any
	)

	stringBuilder.WriteString(`
		INSERT INTO items (
			id, name
		) VALUES
	`)

	argPos := 1
	for i, item := range items {
		if i > 0 {
			stringBuilder.WriteString(",")
		}

		stringBuilder.WriteString(fmt.Sprintf(
			"($%d,$%d)",
			argPos, argPos+1,
		))

		args = append(args,
			item.ID,
			item.Name,
		)

		argPos += 2
	}

	stringBuilder.WriteString(`
		ON CONFLICT (id) DO UPDATE SET
				name = EXCLUDED.name
	`)

	ctx, cancel := context.WithTimeout(ctx, repository.QueryRequestTimeout)
	defer cancel()

	_, err := i.db.ExecContext(ctx, stringBuilder.String(), args...)
	return err
}
