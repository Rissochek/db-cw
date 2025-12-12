package service

import (
	"context"

	"github.com/Rissochek/db-cw/internal/model"
)

type Service struct {
	faker faker
	repo  repo
}

func NewService(faker faker, repo repo) *Service{
	return &Service{
		faker: faker,
		repo: repo,
	}
} 

type faker interface {
	GenerateFakeUsers(toGen int) (users []model.User)
	GenerateFakeListings(toGen int, users []model.User) (listings []model.Listing, listingsMap map[int][]model.Listing)
	GenerateFakeBookings(toGen int, users []model.User, listings []model.Listing, listingsMap map[int][]model.Listing) (bookings []model.Booking)
	GenerateFakeReviews(toGen int, bookings []model.Booking, listings []model.Listing) (reviews []model.Review)
}

type repo interface {
	CreateUsers(ctx context.Context, users []model.User) error
	CreateListings(ctx context.Context, listings []model.Listing) error
	CreateBookings(ctx context.Context, bookings []model.Booking) error
	CreateReviews(ctx context.Context, reviews []model.Review) error
}
