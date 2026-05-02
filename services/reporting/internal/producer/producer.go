package producer

import (
	"context"
	"skinbaron-analyzer/pkg/messaging/jobs"
)

type SyncJobsProducer interface {
	PublishJobRequested(ctx context.Context, jobID string, jobType jobs.SyncJobType) error
}
