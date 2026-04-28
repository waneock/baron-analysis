package repository

import (
	"context"
	"database/sql"
	"skinbaron-analyzer/services/reporting/internal/domain"
	"skinbaron-analyzer/services/reporting/internal/repository/postgres"
	"time"
)

const (
	QueryRequestTimeout = 15 * time.Second
)

type JobsRepository interface {
	Create(ctx context.Context, job domain.SyncJob) error
	GetByID(ctx context.Context, id string) (*domain.SyncJob, error)
}

type Repository struct {
	JobsRepo JobsRepository
}

func New(db *sql.DB) *Repository {
	return &Repository{
		JobsRepo: postgres.NewJobsRepo(db),
	}
}
