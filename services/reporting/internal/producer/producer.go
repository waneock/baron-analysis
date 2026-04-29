package producer

import (
	"context"
	"skinbaron-analyzer/services/reporting/internal/domain"
)

type SyncJobsProducer interface {
	PublishJobRequested(ctx context.Context, jobID string, jobType domain.SyncJobType) error
}
