package domain

import "time"

type SyncJobStatus string

const (
	SyncJobStatusPending SyncJobStatus = "pending"
	SyncJobStatusRunning SyncJobStatus = "running"
	SyncJobStatusDone    SyncJobStatus = "done"
	SyncJobStatusFailed  SyncJobStatus = "failed"
)

type SyncJob struct {
	ID         string
	JobType    string
	Status     SyncJobStatus
	Message    string
	StartedAt  time.Time
	FinishedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
