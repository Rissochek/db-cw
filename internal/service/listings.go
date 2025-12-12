package service

import (
	"context"

	"github.com/Rissochek/db-cw/internal/model"
)

func (s *Service) CreateListing(ctx context.Context, listing *model.Listing) error {
	listing.IsAvailable = true

	return s.repo.CreateListing(ctx, listing)
}

func (s *Service) GetListingByID(ctx context.Context, id int) (*model.Listing, error) {
	return s.repo.GetListingByID(ctx, id)
}

func (s *Service) UpdateListing(ctx context.Context, listing *model.Listing) error {
	dbListing, err := s.repo.GetListingByID(ctx, listing.ID)
	if err != nil {
		return err
	}

	listing.HostID = dbListing.HostID
	listing.Address = dbListing.Address
	
	return s.repo.UpdateListing(ctx, listing)
}

func (s *Service) DeleteListing(ctx context.Context, id int) error {
	return s.repo.DeleteListing(ctx, id)
}

func (s *Service) CreateListings(ctx context.Context, listings []model.Listing) error {
	return s.repo.CreateListings(ctx, listings)
}