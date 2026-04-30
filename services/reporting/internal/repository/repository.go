package repository

import (
	"context"
	"database/sql"
	"skinbaron-analyzer/pkg/messaging/jobs"
	"skinbaron-analyzer/services/reporting/internal/repository/postgres"
	"time"
)

const (
	QueryRequestTimeout = 15 * time.Second
)

type JobsRepository interface {
	Create(ctx context.Context, job jobs.SyncJob) error
	GetByID(ctx context.Context, id string) (*jobs.SyncJob, error)
	MarkFailed(ctx context.Context, id, msg string) error
}

type Repository struct {
	JobsRepo JobsRepository
}

func New(db *sql.DB) *Repository {
	return &Repository{
		JobsRepo: postgres.NewJobsRepo(db),
	}
}
