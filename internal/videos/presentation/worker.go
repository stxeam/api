package presentation

import (
	"context"
	"log/slog"
	"time"

	"go-starter/internal/videos/application"
)

type ProducerWorker struct {
	pollPendingVideos *application.PollPendingVideos
	isRunning         bool
}

func NewProducerWorker(
	pollPendingVideos *application.PollPendingVideos,
) *ProducerWorker {
	return &ProducerWorker{
		pollPendingVideos: pollPendingVideos,
	}
}

func (w *ProducerWorker) Start(ctx context.Context) {
	slog.Info("[Producer] worker started")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("[Producer] worker stopped")
			return
		case <-ticker.C:
			if w.isRunning {
				slog.Info("[Producer] previous cycle still running, skipping tick")
				continue
			}
			w.run(ctx)
		}
	}
}

func (w *ProducerWorker) run(ctx context.Context) {
	w.isRunning = true
	defer func() { w.isRunning = false }()

	result, err := w.pollPendingVideos.Execute(ctx)
	if err != nil {
		slog.Error("[Producer] error during poll cycle", "error", err)
		return
	}

	if len(result.ProcessedIDs) > 0 {
		slog.Info("[Producer] published jobs", "count", len(result.ProcessedIDs), "ids", result.ProcessedIDs)
	}
	if len(result.FailedIDs) > 0 {
		slog.Error("[Producer] failed to publish jobs", "count", len(result.FailedIDs), "ids", result.FailedIDs)
	}
}
