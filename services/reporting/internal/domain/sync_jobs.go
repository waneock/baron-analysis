package domain

import "time"

type SyncJobStatus string

const (
	SyncJobStatusPending SyncJobStatus = "pending"
	SyncJobStatusRunning SyncJobStatus = "running"
	SyncJobStatusDone    SyncJobStatus = "done"
	SyncJobStatusFailed  SyncJobStatus = "failed"
)

type SyncJobType string

const (
	SyncJobTypeSyncOffers    SyncJobType = "sync_offers"
	SyncJobTypeSyncItemSales SyncJobType = "sync_item_sales"
)

type SyncJob struct {
	ID         string
	JobType    SyncJobType
	Status     SyncJobStatus
	Message    string
	StartedAt  time.Time
	FinishedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
