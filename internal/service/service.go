package service

import (
	"context"

	"github.com/Rissochek/db-cw/internal/model"
)

type Service struct {
	faker Faker
	repo  Repo
}

func NewService(faker Faker, repo Repo) *Service {
	return &Service{
		faker: faker,
		repo:  repo,
	}
}

type Faker interface {
	GenerateFakeUsers(toGen int) (users []model.User)
	GenerateFakeListings(toGen int, users []model.User) (listings []model.Listing, listingsMap map[int][]model.Listing)
	GenerateFakeBookings(toGen int, users []model.User, listings []model.Listing, listingsMap map[int][]model.Listing) (bookings []model.Booking)
	GenerateFakeReviews(toGen int, bookings []model.Booking, listings []model.Listing) (reviews []model.Review)
}

type Repo interface {
	CreateUsers(ctx context.Context, users []model.User) error
	CreateListings(ctx context.Context, listings []model.Listing) error
	CreateBookings(ctx context.Context, bookings []model.Booking) error
	CreateReviews(ctx context.Context, reviews []model.Review) error

	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int) error

	CreateListing(ctx context.Context, listing *model.Listing) error
	GetListingByID(ctx context.Context, id int) (*model.Listing, error)
	UpdateListing(ctx context.Context, listing *model.Listing) error
	DeleteListing(ctx context.Context, id int) error

	CreateBooking(ctx context.Context, booking *model.Booking) error
	GetBookingByID(ctx context.Context, bookingID int) (*model.Booking, error)
	UpdateBooking(ctx context.Context, booking *model.Booking) error
	DeleteBooking(ctx context.Context, bookingID int) error
	GetBookingsByID(ctx context.Context, bookingIDs []int) ([]model.Booking, error)
	GetBookingsByListingID(ctx context.Context, listingID int) ([]model.Booking, error)

	CreateReview(ctx context.Context, review *model.Review) error
	GetReviewByID(ctx context.Context, id int) (*model.Review, error)
	UpdateReview(ctx context.Context, review *model.Review) error
	DeleteReview(ctx context.Context, id int) error
}
