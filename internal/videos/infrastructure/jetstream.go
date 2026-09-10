package infrastructure

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

const videoStreamName = "VIDEOS"

func EnsureVideoStream(js jetstream.JetStream) error {
	_, err := js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:     videoStreamName,
		Subjects: []string{"videos.process", "videos.status"},
		MaxAge:   24 * time.Hour,
	})
	return err
}
