package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"skinbaron-analyzer/pkg/messaging/jobs"

	kafkago "github.com/segmentio/kafka-go"
)

type SyncJobRequested struct {
	JobID   string `json:"job_id"`
	JobType string `json:"job_type"`
}

type JobsHandler interface {
	Handle(ctx context.Context, event jobs.SyncJobRequested) error
}

type Consumer struct {
	reader  *kafkago.Reader
	handler JobsHandler
	logger  *slog.Logger
}

func NewConsumer(reader *kafkago.Reader, handler JobsHandler, logger *slog.Logger) *Consumer {
	return &Consumer{
		reader:  reader,
		handler: handler,
		logger:  logger,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	defer c.reader.Close()

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			c.logger.Error("read kafka message",
				"error", err)
			continue
		}

		var event SyncJobRequested
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.logger.Error("unmarshal kafka message",
				"error", err)
			continue
		}

		handlerEvent := syncJobRequestedToHandlerIn(event)

		if err := c.handler.Handle(ctx, handlerEvent); err != nil {
			c.logger.Error("handle kafka message",
				"error", err)
			continue
		}
	}
}

func syncJobRequestedToHandlerIn(input SyncJobRequested) jobs.SyncJobRequested {
	return jobs.SyncJobRequested{
		ID:      input.JobID,
		JobType: jobs.SyncJobType(input.JobType),
	}
}
