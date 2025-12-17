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
	GenerateFakeAmenities() (amenities []model.Amenity)
	GenerateFakeListingAmenities(listings []model.Listing, amenities []model.Amenity) (listingAmenities []model.ListingAmenity)
	GenerateFakeFavorites(toGen int, users []model.User, listings []model.Listing) (favorites []model.Favorite)
	GenerateFakePayments(bookings []model.Booking) (payments []model.Payment)
	GenerateFakeImages(listings []model.Listing) (images []model.Image)
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

	CreateAmenity(ctx context.Context, amenity *model.Amenity) error
	GetAmenityByID(ctx context.Context, id int) (*model.Amenity, error)
	GetAllAmenities(ctx context.Context) ([]model.Amenity, error)
	UpdateAmenity(ctx context.Context, amenity *model.Amenity) error
	DeleteAmenity(ctx context.Context, id int) error
	AddAmenityToListing(ctx context.Context, listingID int, amenityID int) error
	RemoveAmenityFromListing(ctx context.Context, listingID int, amenityID int) error
	GetAmenitiesByListingID(ctx context.Context, listingID int) ([]model.Amenity, error)

	CreateFavorite(ctx context.Context, favorite *model.Favorite) error
	GetFavoriteByID(ctx context.Context, id int) (*model.Favorite, error)
	GetFavoritesByUserID(ctx context.Context, userID int) ([]model.Favorite, error)
	GetFavoriteByUserAndListing(ctx context.Context, userID int, listingID int) (*model.Favorite, error)
	DeleteFavorite(ctx context.Context, id int) error
	DeleteFavoriteByUserAndListing(ctx context.Context, userID int, listingID int) error

	CreateListingAmenities(ctx context.Context, listingAmenities []model.ListingAmenity) error
	CreateFavorites(ctx context.Context, favorites []model.Favorite) error

	CreatePayment(ctx context.Context, payment *model.Payment) error
	GetPaymentByID(ctx context.Context, paymentID int) (*model.Payment, error)
	GetPaymentsByBookingID(ctx context.Context, bookingID int) ([]model.Payment, error)
	UpdatePayment(ctx context.Context, payment *model.Payment) error
	DeletePayment(ctx context.Context, paymentID int) error
	CreatePayments(ctx context.Context, payments []model.Payment) error

	UpdateBookingIsPaid(ctx context.Context, bookingID int, isPaid bool) error

	CreateImage(ctx context.Context, image *model.Image) error
	GetImageByID(ctx context.Context, imageID int) (*model.Image, error)
	GetImagesByListingID(ctx context.Context, listingID int) ([]model.Image, error)
	UpdateImage(ctx context.Context, image *model.Image) error
	DeleteImage(ctx context.Context, imageID int) error
	CreateImages(ctx context.Context, images []model.Image) error
}
