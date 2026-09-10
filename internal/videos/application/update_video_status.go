package application

import (
	"context"

	"go-starter/internal/videos/domain"
)

type UpdateVideoStatusInput struct {
	VideoID   string
	Status    domain.VideoStatus
	Qualities []domain.VideoQuality
}

type UpdateVideoStatus struct {
	videoRepo domain.IVideoRepository
}

func NewUpdateVideoStatus(videoRepo domain.IVideoRepository) *UpdateVideoStatus {
	return &UpdateVideoStatus{videoRepo: videoRepo}
}

func (uc *UpdateVideoStatus) Execute(ctx context.Context, input UpdateVideoStatusInput) error {
	video, err := uc.videoRepo.FindByID(ctx, input.VideoID)
	if err != nil {
		return err
	}
	if video == nil {
		return &domain.VideoNotFoundError{ID: input.VideoID}
	}

	if err := video.TransitionTo(input.Status); err != nil {
		return err
	}
	if len(input.Qualities) > 0 {
		video.SetQualities(input.Qualities)
	}

	_, err = uc.videoRepo.Update(ctx, video)
	return err
}
