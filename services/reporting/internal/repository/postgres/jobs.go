package postgres

import (
	"context"
	"database/sql"
	"skinbaron-analyzer/services/reporting/internal/domain"
	"time"
)

type JobsRepo struct {
	db *sql.DB
}

func NewJobsRepo(db *sql.DB) *JobsRepo {
	return &JobsRepo{
		db: db,
	}
}

func (j *JobsRepo) Create(ctx context.Context, job domain.SyncJob) error {
	query := `
			INSERT INTO sync_jobs (
				id,
				job_type,
				status,
				message,
				started_at,
				finished_at,
				created_at,
				updated_at
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				NULL,
				NULL,
				now(),
				now()
			);
		`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) // TODO replace with a constant
	defer cancel()

	_, err := j.db.ExecContext(
		ctx,
		query,
		job.ID,
		job.JobType,
		job.Status,
		job.Message,
	)

	return err
}

func (j *JobsRepo) GetByID(ctx context.Context, id string) (*domain.SyncJob, error) {
	query := `
		SELECT
			id,
			job_type,
			status,
			message,
			started_at,
			finished_at,
			created_at,
			updated_at
		FROM 
			sync_jobs
		WHERE
			id = $1
	`

	var job domain.SyncJob

	err := j.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&job.ID,
		&job.JobType,
		&job.Status,
		&job.Message,
		&job.StartedAt,
		&job.FinishedAt,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &job, nil
}
