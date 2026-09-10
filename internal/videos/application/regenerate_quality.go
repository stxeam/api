package application

import (
	"context"

	"go-starter/internal/videos/domain"
)

type RegenerateQualityInput struct {
	VideoID string
	Quality string
}

type RegenerateQuality struct {
	videoRepo domain.IVideoRepository
	storage   domain.IStorageAdapter
	publisher domain.IPublisher
}

func NewRegenerateQuality(
	videoRepo domain.IVideoRepository,
	storage domain.IStorageAdapter,
	publisher domain.IPublisher,
) *RegenerateQuality {
	return &RegenerateQuality{
		videoRepo: videoRepo,
		storage:   storage,
		publisher: publisher,
	}
}

func (uc *RegenerateQuality) Execute(ctx context.Context, input RegenerateQualityInput) error {
	video, err := uc.videoRepo.FindByID(ctx, input.VideoID)
	if err != nil {
		return err
	}
	if video == nil {
		return &domain.VideoNotFoundError{ID: input.VideoID}
	}

	rawKey := "raws/" + video.GetID() + ".mp4"
	exists, err := uc.storage.ObjectExists(ctx, rawKey)
	if err != nil {
		return err
	}
	if !exists {
		return &domain.VideoRawNotFoundError{ID: input.VideoID}
	}

	job := domain.VideoProcessingJob{
		VideoID:            video.GetID(),
		RequestedQualities: []string{input.Quality},
		IsAppend:           true,
	}

	return uc.publisher.Publish(ctx, job)
}
