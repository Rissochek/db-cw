package service

import (
	"context"

	"github.com/Rissochek/db-cw/internal/model"
)

func (s *Service) CreateReview(ctx context.Context, review *model.Review) error {
	dbBooking, err := s.repo.GetBookingByID(ctx, review.BookingID)
	if err != nil {
		return err
	}

	review.UserID = dbBooking.GuestID
	return s.repo.CreateReview(ctx, review)
}

func (s *Service) GetReviewByID(ctx context.Context, id int) (*model.Review, error) {
	return s.repo.GetReviewByID(ctx, id)
}

func (s *Service) UpdateReview(ctx context.Context, review *model.Review) error {
	dbReview, err := s.repo.GetReviewByID(ctx, review.ID)
	if err != nil {
		return err
	}

	review.BookingID = dbReview.BookingID
	review.UserID = dbReview.UserID
	return s.repo.UpdateReview(ctx, review)
}

func (s *Service) DeleteReview(ctx context.Context, id int) error {
	return s.repo.DeleteReview(ctx, id)
}
