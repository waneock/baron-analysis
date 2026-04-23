package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"skinbaron-analyzer/services/parsing/internal/repository"
	"strings"
)

type ItemWearsRepo struct {
	db *sql.DB
}

func NewItemWearsRepo(db *sql.DB) *ItemWearsRepo {
	return &ItemWearsRepo{
		db: db,
	}
}

func (w *ItemWearsRepo) CreateMany(ctx context.Context, wears []domain.ItemWearRow) error {
	if len(wears) == 0 {
		return nil
	}

	var (
		stringBuilder strings.Builder
		args          []any
	)

	stringBuilder.WriteString(`
		INSERT INTO item_wears (
			item_id, name
		) VALUES
	`)

	argPos := 1
	for i, wear := range wears {
		if i > 0 {
			stringBuilder.WriteString(",")
		}

		stringBuilder.WriteString(fmt.Sprintf(
			"($%d,$%d)",
			argPos, argPos+1,
		))

		args = append(args,
			wear.ID,
			wear.Name,
		)

		argPos += 2
	}

	stringBuilder.WriteString(`
		ON CONFLICT (id) DO UPDATE SET
				name = EXCLUDED.name
	`)

	ctx, cancel := context.WithTimeout(ctx, repository.QueryRequestTimeout)
	defer cancel()

	_, err := w.db.ExecContext(ctx, stringBuilder.String(), args...)
	return err
}
