package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"

	"go-starter/internal/videos/domain"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsPublisher struct {
	js jetstream.JetStream
}

func NewNatsPublisher(natsURL string) (*NatsPublisher, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("create jetstream: %w", err)
	}
	if err := EnsureVideoStream(js); err != nil {
		return nil, fmt.Errorf("ensure video stream: %w", err)
	}

	return &NatsPublisher{js: js}, nil
}

func (p *NatsPublisher) Publish(ctx context.Context, job domain.VideoProcessingJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal job: %w", err)
	}

	_, err = p.js.Publish(ctx, "videos.process", data)
	if err != nil {
		return fmt.Errorf("publish job: %w", err)
	}

	return nil
}
