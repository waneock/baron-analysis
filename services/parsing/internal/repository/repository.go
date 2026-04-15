package repository

import (
	"context"
	"database/sql"
	"skinbaron-analyzer/services/parsing/internal/domain"
	"skinbaron-analyzer/services/parsing/internal/repository/postgres"
	"time"
)

const (
	QueryRequestTimeout = 5 * time.Second
)

type OffersRepository interface {
	CreateMany(ctx context.Context, offers []domain.Offer) error
}

type Repo struct {
	OffersRepository OffersRepository
}

func New(db *sql.DB) *Repo {
	return &Repo{
		OffersRepository: postgres.NewOffersRepo(db),
	}
}
