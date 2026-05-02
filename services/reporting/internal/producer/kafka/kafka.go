package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"skinbaron-analyzer/pkg/messaging/jobs"

	kafkago "github.com/segmentio/kafka-go"
)

const (
	TopicSyncJobsRequested = "sync.jobs.requested"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(writer *kafkago.Writer) *Producer {
	return &Producer{
		writer: writer,
	}
}

func (p *Producer) PublishJobRequested(ctx context.Context, jobID string, jobType jobs.SyncJobType) error {
	msg := SyncJobRequested{
		JobID:   jobID,
		JobType: string(jobType),
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal kafka message: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafkago.Message{
		Topic: TopicSyncJobsRequested,
		Key:   []byte(jobID),
		Value: payload,
	})
	if err != nil {
		return fmt.Errorf("write kafka message: %w", err)
	}

	return nil
}
