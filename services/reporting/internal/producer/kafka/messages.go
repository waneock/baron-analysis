package kafka

type SyncJobRequested struct {
	JobID   string `json:"job_id"`
	JobType string `json:"job_type"`
}
