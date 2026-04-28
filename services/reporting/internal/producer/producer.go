package producer

import "context"

type SyncJobsProducer interface {
	PublishJobRequested(ctx context.Context, jobID, jobType string) error
}
