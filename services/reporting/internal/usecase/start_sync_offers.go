package usecase

import (
	"context"
	"skinbaron-analyzer/services/reporting/internal/domain"

	"github.com/google/uuid"
)

const (
	MsgFailToPublishJob = "failed to publish kafka message"
)

type SyncJobsRepository interface {
	Create(ctx context.Context, job domain.SyncJob) error
	MarkFailed(ctx context.Context, jobID, msg string) error
}

type SyncJobsProducer interface {
	PublishJobRequested(ctx context.Context, jobID string, jobType domain.SyncJobType) error
}

type StartSyncJob struct {
	repo     SyncJobsRepository
	producer SyncJobsProducer
}

func NewSyncOffers(repo SyncJobsRepository, producer SyncJobsProducer) *StartSyncJob {
	return &StartSyncJob{
		repo:     repo,
		producer: producer,
	}
}

func (uc *StartSyncJob) Execute(ctx context.Context, jobType domain.SyncJobType) (string, error) {
	jobID := uuid.NewString()

	job := domain.SyncJob{
		ID:      jobID,
		JobType: jobType,
		Status:  domain.SyncJobStatusPending,
	}

	if err := uc.repo.Create(ctx, job); err != nil {
		return "", err
	}

	if err := uc.producer.PublishJobRequested(ctx, jobID, jobType); err != nil {
		_ = uc.repo.MarkFailed(ctx, jobID, MsgFailToPublishJob)
		return "", err
	}

	return jobID, nil
}
