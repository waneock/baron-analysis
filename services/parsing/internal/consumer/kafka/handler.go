package kafka

import (
	"context"
	"skinbaron-analyzer/pkg/messaging/jobs"
)

type SyncOffersRunner interface {
	Execute(ctx context.Context, jobID string)
}

type SyncItemsRunner interface {
	Execute(ctx context.Context, jobID string)
}

type SyncItemSalesRunner interface {
	Execute(ctx context.Context, jobID string)
}

type JobEventHandler struct {
	syncOffers    SyncOffersRunner
	syncItems     SyncItemsRunner
	syncItemSales SyncItemSalesRunner
}

func NewJobsEventHandler(syncOffers SyncOffersRunner, syncItems SyncItemsRunner, syncItemSales SyncItemSalesRunner) *JobEventHandler {
	return &JobEventHandler{
		syncOffers:    syncOffers,
		syncItems:     syncItems,
		syncItemSales: syncItemSales,
	}
}

func (h *JobEventHandler) Handle(ctx context.Context, event jobs.SyncJobRequested) {
	switch event.JobType {
	case jobs.SyncJobTypeSyncOffers:
		h.syncOffers.Execute(ctx, event.ID)
	case jobs.SyncJobTypeSyncItems:
		h.syncItems.Execute(ctx, event.ID)
	case jobs.SyncJobTypeSyncItemSales:
		h.syncItemSales.Execute(ctx, event.ID)
	}
}
