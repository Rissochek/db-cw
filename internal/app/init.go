package app

import (
	"context"
	"time"

	"github.com/Rissochek/db-cw/api/handler"
	_ "github.com/Rissochek/db-cw/docs"
	"github.com/Rissochek/db-cw/internal/faking"
	"github.com/Rissochek/db-cw/internal/repository/postgres"
	"github.com/Rissochek/db-cw/internal/service"
	"github.com/Rissochek/db-cw/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var (
	seed           = int64(42)
	isGen          = false
	migrationsPath = "./internal/migrations"
)

func InitApp() *App {
	utils.LoadEnvFile()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	faker := faking.NewDataFaker(seed)

	dsn := postgres.CreateDsnFromEnv()
	conn := postgres.CreateConnection(dsn)
	repo := postgres.NewPostgres(conn)
	repo.RunMigrations(migrationsPath)

	service := service.NewService(faker, repo)

	if isGen {
		service.FillDatabase(ctx, seed)
	}

	handler := handler.NewHandler(service)

	app := NewApp(handler)

	return app
}

func (app *App) SetupRoutes() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api")

	api.POST("/users", app.handler.CreateUser)
	api.POST("/users/batch", app.handler.BatchImportUsers)
	api.GET("/users/:id", app.handler.GetUserByID)
	api.PUT("/users/:id", app.handler.UpdateUser)
	api.DELETE("/users/:id", app.handler.DeleteUser)

	api.POST("/listings", app.handler.CreateListing)
	api.POST("/listings/batch", app.handler.BatchImportListings)
	api.GET("/listings/:id", app.handler.GetListingByID)
	api.PUT("/listings/:id", app.handler.UpdateListing)
	api.DELETE("/listings/:id", app.handler.DeleteListing)

	api.POST("/bookings", app.handler.CreateBooking)
	api.POST("/bookings/batch", app.handler.BatchImportBookings)
	api.GET("/bookings/:id", app.handler.GetBookingByID)
	api.PUT("/bookings/:id", app.handler.UpdateBooking)
	api.DELETE("/bookings/:id", app.handler.DeleteBooking)

	api.POST("/reviews", app.handler.CreateReview)
	api.POST("/reviews/batch", app.handler.BatchImportReviews)
	api.GET("/reviews/:id", app.handler.GetReviewByID)
	api.PUT("/reviews/:id", app.handler.UpdateReview)
	api.DELETE("/reviews/:id", app.handler.DeleteReview)

	api.POST("/amenities", app.handler.CreateAmenity)
	api.GET("/amenities", app.handler.GetAllAmenities)
	api.GET("/amenities/:id", app.handler.GetAmenityByID)
	api.PUT("/amenities/:id", app.handler.UpdateAmenity)
	api.DELETE("/amenities/:id", app.handler.DeleteAmenity)
	api.POST("/listings/:listing_id/amenities/:amenity_id", app.handler.AddAmenityToListing)
	api.DELETE("/listings/:listing_id/amenities/:amenity_id", app.handler.RemoveAmenityFromListing)
	api.GET("/listings/:listing_id/amenities", app.handler.GetAmenitiesByListingID)

	api.POST("/favorites", app.handler.CreateFavorite)
	api.GET("/favorites/:id", app.handler.GetFavoriteByID)
	api.GET("/users/:user_id/favorites", app.handler.GetFavoritesByUserID)
	api.DELETE("/favorites/:id", app.handler.DeleteFavorite)
	api.DELETE("/users/:user_id/favorites/:listing_id", app.handler.DeleteFavoriteByUserAndListing)

	api.POST("/payments", app.handler.CreatePayment)
	api.GET("/payments/:id", app.handler.GetPaymentByID)
	api.GET("/bookings/:booking_id/payments", app.handler.GetPaymentsByBookingID)
	api.PUT("/payments/:id", app.handler.UpdatePayment)
	api.DELETE("/payments/:id", app.handler.DeletePayment)

	api.POST("/images", app.handler.CreateImage)
	api.GET("/images/:id", app.handler.GetImageByID)
	api.GET("/listings/:listing_id/images", app.handler.GetImagesByListingID)
	api.PUT("/images/:id", app.handler.UpdateImage)
	api.DELETE("/images/:id", app.handler.DeleteImage)

	api.GET("/functions/hosts/:host_id/revenue", app.handler.GetHostTotalRevenue)
	api.GET("/functions/guests/:guest_id/total-spent", app.handler.GetGuestTotalSpent)
	api.GET("/functions/hosts/:host_id/average-rating", app.handler.GetHostAverageRating)
	api.GET("/functions/listings/:listing_id/active-bookings", app.handler.GetListingActiveBookingsCount)

	api.GET("/reports/listings-statistics", app.handler.GetListingsStatisticsReport)
	api.GET("/reports/hosts-performance", app.handler.GetHostsPerformanceReport)
	api.GET("/reports/bookings", app.handler.GetBookingsReport)
	api.GET("/reports/payments-summary", app.handler.GetPaymentsSummaryReport)

	api.POST("/procedures/create-booking-with-payment", app.handler.CreateBookingWithPayment)
	api.POST("/procedures/payments/:id/confirm", app.handler.ConfirmPayment)
	api.POST("/procedures/bookings/:id/cancel-with-refund", app.handler.CancelBookingWithRefund)

	return e
}
