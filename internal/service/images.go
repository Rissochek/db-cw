package service

import (
	"context"
	"time"

	"github.com/Rissochek/db-cw/internal/model"
)

func (s *Service) CreateImage(ctx context.Context, image *model.Image) error {
	_, err := s.repo.GetListingByID(ctx, image.ListingID)
	if err != nil {
		return err
	}

	if image.UploadedAt.IsZero() {
		image.UploadedAt = time.Now()
	}

	return s.repo.CreateImage(ctx, image)
}

func (s *Service) GetImageByID(ctx context.Context, imageID int) (*model.Image, error) {
	return s.repo.GetImageByID(ctx, imageID)
}

func (s *Service) GetImagesByListingID(ctx context.Context, listingID int) ([]model.Image, error) {
	return s.repo.GetImagesByListingID(ctx, listingID)
}

func (s *Service) UpdateImage(ctx context.Context, image *model.Image) error {
	dbImage, err := s.repo.GetImageByID(ctx, image.ImageID)
	if err != nil {
		return err
	}

	image.ListingID = dbImage.ListingID
	if image.UploadedAt.IsZero() {
		image.UploadedAt = dbImage.UploadedAt
	}

	return s.repo.UpdateImage(ctx, image)
}

func (s *Service) DeleteImage(ctx context.Context, imageID int) error {
	return s.repo.DeleteImage(ctx, imageID)
}

func (s *Service) CreateImages(ctx context.Context, images []model.Image) error {
	return s.repo.CreateImages(ctx, images)
}
