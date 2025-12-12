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
	seed = int64(42)
)

func InitApp() *App {
	utils.LoadEnvFile()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	faker := faking.NewDataFaker(seed)

	dsn := postgres.CreateDsnFromEnv()
	conn := postgres.CreateConnection(dsn)
	repo := postgres.NewPostgres(conn)

	service := service.NewService(faker, repo)
	service.FillDatabase(ctx, seed)

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

	return e
}
