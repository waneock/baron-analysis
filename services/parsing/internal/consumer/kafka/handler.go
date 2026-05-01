package kafka

import (
	"context"
	"skinbaron-analyzer/pkg/messaging/jobs"
)

type SyncOffersRunner interface {
	Execute(ctx context.Context, jobID string)
}

type JobEventHandler struct {
	syncOffers SyncOffersRunner
}

func NewJobsEventHandler(syncOffers SyncOffersRunner) *JobEventHandler {
	return &JobEventHandler{
		syncOffers: syncOffers,
	}
}

func (h *JobEventHandler) Handle(ctx context.Context, event jobs.SyncJobRequested) {
	switch event.JobType {
	case "sync_offers":
		h.syncOffers.Execute(ctx, event.ID)
	}
}
