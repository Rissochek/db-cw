package service

import (
	"context"

	"github.com/Rissochek/db-cw/internal/model"
)

func (s *Service) CreateAmenity(ctx context.Context, amenity *model.Amenity) error {
	return s.repo.CreateAmenity(ctx, amenity)
}

func (s *Service) GetAmenityByID(ctx context.Context, id int) (*model.Amenity, error) {
	return s.repo.GetAmenityByID(ctx, id)
}

func (s *Service) GetAllAmenities(ctx context.Context) ([]model.Amenity, error) {
	return s.repo.GetAllAmenities(ctx)
}

func (s *Service) UpdateAmenity(ctx context.Context, amenity *model.Amenity) error {
	return s.repo.UpdateAmenity(ctx, amenity)
}

func (s *Service) DeleteAmenity(ctx context.Context, id int) error {
	return s.repo.DeleteAmenity(ctx, id)
}

func (s *Service) AddAmenityToListing(ctx context.Context, listingID int, amenityID int) error {
	_, err := s.repo.GetListingByID(ctx, listingID)
	if err != nil {
		return err
	}

	_, err = s.repo.GetAmenityByID(ctx, amenityID)
	if err != nil {
		return err
	}

	return s.repo.AddAmenityToListing(ctx, listingID, amenityID)
}

func (s *Service) RemoveAmenityFromListing(ctx context.Context, listingID int, amenityID int) error {
	return s.repo.RemoveAmenityFromListing(ctx, listingID, amenityID)
}

func (s *Service) GetAmenitiesByListingID(ctx context.Context, listingID int) ([]model.Amenity, error) {
	return s.repo.GetAmenitiesByListingID(ctx, listingID)
}

func (s *Service) CreateListingAmenities(ctx context.Context, listingAmenities []model.ListingAmenity) error {
	return s.repo.CreateListingAmenities(ctx, listingAmenities)
}
