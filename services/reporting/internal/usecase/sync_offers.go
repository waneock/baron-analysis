package usecase

import (
	"context"
	"skinbaron-analyzer/pkg/messaging/jobs"

	"github.com/google/uuid"
)

const (
	MsgFailToPublishJob = "failed to publish kafka message"
)

type SyncJobsRepository interface {
	Create(ctx context.Context, job jobs.SyncJob) error
	MarkFailed(ctx context.Context, jobID, msg string) error
}

type SyncJobsProducer interface {
	PublishJobRequested(ctx context.Context, jobID string, jobType jobs.SyncJobType) error
}

type SyncOffersJob struct {
	repo     SyncJobsRepository
	producer SyncJobsProducer
}

//TODO: review this file, make a generic service for working with syncs, since the logic is the same

func NewSyncOffers(repo SyncJobsRepository, producer SyncJobsProducer) *SyncOffersJob {
	return &SyncOffersJob{
		repo:     repo,
		producer: producer,
	}
}

func (uc *SyncOffersJob) Execute(ctx context.Context, jobType jobs.SyncJobType) (string, error) {
	jobID := uuid.NewString()

	job := jobs.SyncJob{
		ID:      jobID,
		JobType: jobType,
		Status:  jobs.SyncJobStatusPending,
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
