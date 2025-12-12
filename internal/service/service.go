package service

import "github.com/Rissochek/db-cw/internal/model"

type Service struct {
	faker Faker
}

type Faker interface {
	GenerateFakeUsers(toGen int) (users []model.User)
	GenerateFakeListings(toGen int, users []model.User) (listings []model.Listing, listingsMap map[int][]model.Listing)
	GenerateFakeBookings(toGen int, users []model.User, listings []model.Listing, listingsMap map[int][]model.Listing) (bookings []model.Booking)
	GenerateFakeReviews(toGen int, bookings []model.Booking, listings []model.Listing) (reviews []model.Review)
}
