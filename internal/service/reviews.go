package service

import (
	"context"
	"fmt"

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

func (s *Service) CreateReviews(ctx context.Context, reviews []model.Review) error {
	bookingIDsMap := make(map[int]bool)
	bookingIDs := make([]int, 0, len(reviews))

	for i := range reviews {
		if !bookingIDsMap[reviews[i].BookingID] {
			bookingIDsMap[reviews[i].BookingID] = true
			bookingIDs = append(bookingIDs, reviews[i].BookingID)
		}
	}

	bookings, err := s.repo.GetBookingsByID(ctx, bookingIDs)
	if err != nil {
		return err
	}

	bookingsMap := make(map[int]*model.Booking, len(bookings))
	for i := range bookings {
		bookingsMap[bookings[i].BookingID] = &bookings[i]
	}

	for i := range reviews {
		booking, ok := bookingsMap[reviews[i].BookingID]
		if !ok {
			return fmt.Errorf("booking with id %d not found", reviews[i].BookingID)
		}
		reviews[i].UserID = booking.GuestID
	}

	return s.repo.CreateReviews(ctx, reviews)
}
