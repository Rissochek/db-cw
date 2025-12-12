package handler

import (
	"context"

	"github.com/Rissochek/db-cw/internal/model"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

type Service interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int) error
	CreateUsers(ctx context.Context, users []model.User) error

	CreateListing(ctx context.Context, listing *model.Listing) error
	GetListingByID(ctx context.Context, id int) (*model.Listing, error)
	UpdateListing(ctx context.Context, listing *model.Listing) error
	DeleteListing(ctx context.Context, id int) error
	CreateListings(ctx context.Context, listings []model.Listing) error

	CreateBooking(ctx context.Context, booking *model.Booking) error
	GetBookingByID(ctx context.Context, bookingID int) (*model.Booking, error)
	UpdateBooking(ctx context.Context, booking *model.Booking) error
	DeleteBooking(ctx context.Context, bookingID int) error
	CreateBookings(ctx context.Context, bookings []model.Booking) error

	CreateReview(ctx context.Context, review *model.Review) error
	GetReviewByID(ctx context.Context, id int) (*model.Review, error)
	UpdateReview(ctx context.Context, review *model.Review) error
	DeleteReview(ctx context.Context, id int) error
	CreateReviews(ctx context.Context, reviews []model.Review) error
}
