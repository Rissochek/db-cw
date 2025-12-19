package app

import (
	"github.com/labstack/echo/v4"
)

type App struct {
	handler Handler
}

func NewApp(handler Handler) *App {
	return &App{
		handler: handler,
	}
}

type Handler interface {
	CreateUser(c echo.Context) error
	BatchImportUsers(c echo.Context) error
	GetUserByID(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error

	CreateListing(c echo.Context) error
	BatchImportListings(c echo.Context) error
	GetListingByID(c echo.Context) error
	UpdateListing(c echo.Context) error
	DeleteListing(c echo.Context) error

	CreateBooking(c echo.Context) error
	BatchImportBookings(c echo.Context) error
	GetBookingByID(c echo.Context) error
	UpdateBooking(c echo.Context) error
	DeleteBooking(c echo.Context) error

	CreateReview(c echo.Context) error
	BatchImportReviews(c echo.Context) error
	GetReviewByID(c echo.Context) error
	UpdateReview(c echo.Context) error
	DeleteReview(c echo.Context) error

	CreateAmenity(c echo.Context) error
	GetAmenityByID(c echo.Context) error
	GetAllAmenities(c echo.Context) error
	UpdateAmenity(c echo.Context) error
	DeleteAmenity(c echo.Context) error
	AddAmenityToListing(c echo.Context) error
	RemoveAmenityFromListing(c echo.Context) error
	GetAmenitiesByListingID(c echo.Context) error

	CreateFavorite(c echo.Context) error
	GetFavoriteByID(c echo.Context) error
	GetFavoritesByUserID(c echo.Context) error
	DeleteFavorite(c echo.Context) error
	DeleteFavoriteByUserAndListing(c echo.Context) error

	CreatePayment(c echo.Context) error
	GetPaymentByID(c echo.Context) error
	GetPaymentsByBookingID(c echo.Context) error
	UpdatePayment(c echo.Context) error
	DeletePayment(c echo.Context) error

	CreateImage(c echo.Context) error
	GetImageByID(c echo.Context) error
	GetImagesByListingID(c echo.Context) error
	UpdateImage(c echo.Context) error
	DeleteImage(c echo.Context) error

	GetHostTotalRevenue(c echo.Context) error
	GetGuestTotalSpent(c echo.Context) error
	GetHostAverageRating(c echo.Context) error
	GetListingActiveBookingsCount(c echo.Context) error

	GetListingsStatisticsReport(c echo.Context) error
	GetHostsPerformanceReport(c echo.Context) error
	GetBookingsReport(c echo.Context) error
	GetPaymentsSummaryReport(c echo.Context) error

	CreateBookingWithPayment(c echo.Context) error
	ConfirmPayment(c echo.Context) error
	CancelBookingWithRefund(c echo.Context) error
}
