package postgres

import (
	"context"
	"database/sql"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"skinbaron-analyzer/services/parsing/internal/repository"
)

type MarketSyncSourceRepo struct {
	db *sql.DB
}

func NewMarketSyncSourceRepo(db *sql.DB) *MarketSyncSourceRepo {
	return &MarketSyncSourceRepo{
		db: db,
	}
}

func (m *MarketSyncSourceRepo) Count(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM items i
		JOIN item_wears iw ON iw.item_id = i.id;
		`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryRequestTimeout)
	defer cancel()

	var total int
	err := m.db.QueryRowContext(ctx, query).Scan(&total)
	return total, err
}

func (m *MarketSyncSourceRepo) List(ctx context.Context, limit, offset int) (*[]domain.ItemWearSource, error) {
	query := `
			SELECT
				iw.id,
				i.name,
				iw.name
			FROM items i
			JOIN item_wears iw ON iw.item_id = i.id
			ORDER BY i.name, iw.name
			LIMIT $1 OFFSET $2;
		`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryRequestTimeout)
	defer cancel()

	rows, err := m.db.QueryContext(
		ctx,
		query,
		limit,
		offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.ItemWearSource

	for rows.Next() {
		var item domain.ItemWearSource

		err := rows.Scan(
			&item.WearID,
			&item.Name,
			&item.Wear,
		)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &items, nil
}
