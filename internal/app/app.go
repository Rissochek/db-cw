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
}
