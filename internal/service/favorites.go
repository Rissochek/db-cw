package service

import (
	"context"
	"fmt"

	"github.com/Rissochek/db-cw/internal/model"
)

func (s *Service) CreateFavorite(ctx context.Context, favorite *model.Favorite) error {
	_, err := s.repo.GetUserByID(ctx, favorite.UserID)
	if err != nil {
		return err
	}

	_, err = s.repo.GetListingByID(ctx, favorite.ListingID)
	if err != nil {
		return err
	}

	existing, err := s.repo.GetFavoriteByUserAndListing(ctx, favorite.UserID, favorite.ListingID)
	if err != nil {
		return err
	}

	if existing != nil {
		return fmt.Errorf("favorite already exists")
	}

	return s.repo.CreateFavorite(ctx, favorite)
}

func (s *Service) GetFavoriteByID(ctx context.Context, id int) (*model.Favorite, error) {
	return s.repo.GetFavoriteByID(ctx, id)
}

func (s *Service) GetFavoritesByUserID(ctx context.Context, userID int) ([]model.Favorite, error) {
	return s.repo.GetFavoritesByUserID(ctx, userID)
}

func (s *Service) DeleteFavorite(ctx context.Context, id int) error {
	return s.repo.DeleteFavorite(ctx, id)
}

func (s *Service) DeleteFavoriteByUserAndListing(ctx context.Context, userID int, listingID int) error {
	return s.repo.DeleteFavoriteByUserAndListing(ctx, userID, listingID)
}

func (s *Service) CreateFavorites(ctx context.Context, favorites []model.Favorite) error {
	return s.repo.CreateFavorites(ctx, favorites)
}
