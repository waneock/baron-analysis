package postgres

import (
	"context"
	"database/sql"
	"skinbaron-analyzer/pkg/messaging/jobs"
	"skinbaron-analyzer/services/parsing/internal/repository"
	"time"
)

type JobRepo struct {
	db *sql.DB
}

func NewJobRepo(db *sql.DB) *JobRepo {
	return &JobRepo{
		db: db,
	}
}

func (j *JobRepo) UpdateStatus(ctx context.Context, id string, status jobs.SyncJobStatus) error {
	ctx, cancel := context.WithTimeout(ctx, repository.QueryRequestTimeout)
	defer cancel()

	job, err := j.fetchByID(ctx, id)
	if err != nil {
		return err
	}

	job.Status = status

	query := `
		UPDATE
			sync_jobs
		SET
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
			$5,
			$6,
			$7,
			now()
		);
	`

	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = j.db.ExecContext(
		ctx,
		query,
		job.ID,
		job.JobType,
		job.Status,
		job.Message,
		job.StartedAt,
		job.FinishedAt,
		job.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (j *JobRepo) fetchByID(ctx context.Context, id string) (*jobs.SyncJob, error) {
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

	var job jobs.SyncJob

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
