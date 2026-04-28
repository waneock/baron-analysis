package domain

import "time"

type SyncJob struct {
	ID         string
	JobType    string
	Status     string
	Message    string
	StartedAt  time.Time
	FinishedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
