package service

import (
	"context"
	"time"

	"github.com/Rissochek/db-cw/internal/model"
)

func (s *Service) CreateBooking(ctx context.Context, booking *model.Booking) error {
	dbListing, err := s.repo.GetListingByID(ctx, booking.ListingID)
	if err != nil {
		return err
	}

	booking.HostID = dbListing.HostID
	booking.TotalPrice = countTotalPrice(booking.InDate, booking.OutDate, dbListing.PricePerNight)
	booking.IsPaid = false
	return s.repo.CreateBooking(ctx, booking)
}

func (s *Service) GetBookingByID(ctx context.Context, bookingID int) (*model.Booking, error) {
	return s.repo.GetBookingByID(ctx, bookingID)
}

func (s *Service) UpdateBooking(ctx context.Context, booking *model.Booking) error {
	dbBooking, err := s.repo.GetBookingByID(ctx, booking.BookingID)
	if err != nil {
		return err
	}

	dbListing, err := s.repo.GetListingByID(ctx, dbBooking.ListingID)
	if err != nil {
		return err
	}

	booking.GuestID = dbBooking.GuestID
	booking.HostID = dbBooking.HostID
	booking.ListingID = dbBooking.ListingID
	booking.TotalPrice = countTotalPrice(booking.InDate, booking.OutDate, dbListing.PricePerNight)

	return s.repo.UpdateBooking(ctx, booking)
}

func (s *Service) DeleteBooking(ctx context.Context, bookingID int) error {
	return s.repo.DeleteBooking(ctx, bookingID)
}

func countTotalPrice(inDate, outDate time.Time, pricePerNight float64) float64 {
	diff := outDate.Sub(inDate)

	nights := int(diff.Hours() / 24)
	nights = max(nights, 1)

	totalPrice := float64(nights) * pricePerNight

	return totalPrice
}
