package handler

import (
	"context"
	"time"

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
	DeleteFavorite(ctx context.Context, id int) error
	DeleteFavoriteByUserAndListing(ctx context.Context, userID int, listingID int) error

	CreatePayment(ctx context.Context, payment *model.Payment) error
	GetPaymentByID(ctx context.Context, paymentID int) (*model.Payment, error)
	GetPaymentsByBookingID(ctx context.Context, bookingID int) ([]model.Payment, error)
	UpdatePayment(ctx context.Context, payment *model.Payment) error
	DeletePayment(ctx context.Context, paymentID int) error

	CreateImage(ctx context.Context, image *model.Image) error
	GetImageByID(ctx context.Context, imageID int) (*model.Image, error)
	GetImagesByListingID(ctx context.Context, listingID int) ([]model.Image, error)
	UpdateImage(ctx context.Context, image *model.Image) error
	DeleteImage(ctx context.Context, imageID int) error

	GetHostTotalRevenue(ctx context.Context, hostID int) (float64, error)
	GetGuestTotalSpent(ctx context.Context, guestID int) (float64, error)
	GetHostAverageRating(ctx context.Context, hostID int) (float64, error)
	GetListingActiveBookingsCount(ctx context.Context, listingID int) (int, error)

	GetListingsStatisticsReport(ctx context.Context) ([]model.ListingStatisticsReport, error)
	GetHostsPerformanceReport(ctx context.Context) ([]model.HostPerformanceReport, error)
	GetBookingsReport(ctx context.Context, startDate, endDate *time.Time) ([]model.BookingReport, error)
	GetPaymentsSummaryReport(ctx context.Context, startDate, endDate *time.Time) ([]model.PaymentSummaryReport, error)

	CreateBookingWithPayment(ctx context.Context, listingID, guestID int, inDate, outDate time.Time, paymentMethod string) (*model.CreateBookingWithPaymentResult, error)
	ConfirmPayment(ctx context.Context, paymentID int, transactionID *string) error
	CancelBookingWithRefund(ctx context.Context, bookingID int) error
}
