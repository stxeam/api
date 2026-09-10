package presentation

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"go-starter/internal/videos/application"
	"go-starter/internal/videos/domain"
	videosinfra "go-starter/internal/videos/infrastructure"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type VideoStatusUpdate struct {
	VideoID   string   `json:"video_id"`
	Status    string   `json:"status"`
	Qualities []string `json:"qualities"`
}

type StatusSubscriber struct {
	js                jetstream.JetStream
	updateVideoStatus *application.UpdateVideoStatus
}

func NewStatusSubscriber(natsURL string, repo domain.IVideoRepository) (*StatusSubscriber, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	if err := videosinfra.EnsureVideoStream(js); err != nil {
		return nil, err
	}

	return &StatusSubscriber{
		js:                js,
		updateVideoStatus: application.NewUpdateVideoStatus(repo),
	}, nil
}

func (s *StatusSubscriber) Start(ctx context.Context) error {
	consumer, err := s.js.CreateOrUpdateConsumer(ctx, "VIDEOS", jetstream.ConsumerConfig{
		Durable:       "api-status-subscriber",
		FilterSubject: "videos.status",
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxAckPending: 1,
	})
	if err != nil {
		return err
	}

	consumeContext, err := consumer.Consume(func(msg jetstream.Msg) {
		var update VideoStatusUpdate
		if err := json.Unmarshal(msg.Data(), &update); err != nil {
			slog.Error("failed to decode status update", "err", err)
			_ = msg.Ack()
			return
		}

		targetStatus := domain.VideoStatus(update.Status)
		qualities := make([]domain.VideoQuality, 0, len(update.Qualities))
		for _, quality := range update.Qualities {
			qualities = append(qualities, domain.VideoQuality(quality))
		}
		err := s.updateVideoStatus.Execute(ctx, application.UpdateVideoStatusInput{
			VideoID:   update.VideoID,
			Status:    targetStatus,
			Qualities: qualities,
		})
		if err != nil {
			var notFound *domain.VideoNotFoundError
			if errors.As(err, &notFound) {
				slog.Error("video not found for status update", "id", update.VideoID)
				_ = msg.Ack()
				return
			}

			var invalidTransition *domain.VideoInvalidStatusTransitionError
			if errors.As(err, &invalidTransition) {
				slog.Error("invalid status transition", "id", update.VideoID, "target", targetStatus)
				_ = msg.Ack()
				return
			}

			slog.Error("failed to update video status in db", "err", err)
			_ = msg.Ack()
			return
		}

		if err := msg.Ack(); err != nil {
			slog.Error("failed to acknowledge status update", "err", err)
			return
		}
		slog.Info("updated video status", "id", update.VideoID, "status", targetStatus)
	})
	if err != nil {
		return err
	}
	defer consumeContext.Stop()

	slog.Info("[StatusSubscriber] worker started")
	<-ctx.Done()
	return ctx.Err()
}
