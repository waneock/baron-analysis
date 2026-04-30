package kafka

import (
	"context"
	"skinbaron-analyzer/pkg/messaging/jobs"
)

type SyncOffersRunner interface {
	Execute(ctx context.Context, jobID string) error
}

type JobEventHandler struct {
	syncOffers SyncOffersRunner
}

func NewJobsEventHandler(syncOffers SyncOffersRunner) *JobEventHandler {
	return &JobEventHandler{
		syncOffers: syncOffers,
	}
}

func (h *JobEventHandler) Handle(ctx context.Context, event jobs.SyncJobRequested) error {
	switch event.JobType {
	case "sync_offers":
		return h.syncOffers.Execute(ctx, event.ID)
	default:
		return nil
	}
}
