package application

import (
	"context"

	"go-starter/internal/videos/domain"
)

type TriggerProcessingInput struct {
	VideoID            string
	RequestedQualities []string
}

type TriggerProcessing struct {
	videoRepo domain.IVideoRepository
	publisher domain.IPublisher
}

func NewTriggerProcessing(
	videoRepo domain.IVideoRepository,
	publisher domain.IPublisher,
) *TriggerProcessing {
	return &TriggerProcessing{
		videoRepo: videoRepo,
		publisher: publisher,
	}
}

func (uc *TriggerProcessing) Execute(ctx context.Context, input TriggerProcessingInput) error {
	video, err := uc.videoRepo.FindByID(ctx, input.VideoID)
	if err != nil {
		return err
	}
	if video == nil {
		return &domain.VideoNotFoundError{ID: input.VideoID}
	}

	if err := video.TransitionTo(domain.StatusProcessing); err != nil {
		return err
	}

	_, err = uc.videoRepo.Update(ctx, video)
	if err != nil {
		return err
	}

	job := domain.VideoProcessingJob{
		VideoID:            video.GetID(),
		RequestedQualities: input.RequestedQualities,
		IsAppend:           false,
	}

	return uc.publisher.Publish(ctx, job)
}
