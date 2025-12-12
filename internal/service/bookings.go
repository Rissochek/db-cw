package service

import (
	"context"
	"errors"
	"time"

	"github.com/Rissochek/db-cw/internal/model"
	"go.uber.org/zap"
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

func (s *Service) CreateBookings(ctx context.Context, bookings []model.Booking) error {
	batchBookingsMap := make(map[int][]model.Booking)
	dbBookingsCache := make(map[int][]model.Booking)
	listingsCache := make(map[int]*model.Listing)

	for i := range bookings {
		listingID := bookings[i].ListingID

		dbListing, ok := listingsCache[listingID]
		if !ok {
			var err error
			dbListing, err = s.repo.GetListingByID(ctx, listingID)
			if err != nil {
				return err
			}
			listingsCache[listingID] = dbListing
		}

		dbBookings, ok := dbBookingsCache[listingID]
		if !ok {
			var err error
			dbBookings, err = s.repo.GetBookingsByListingID(ctx, listingID)
			if err != nil {
				return err
			}
			dbBookingsCache[listingID] = dbBookings
		}

		batchBookings := batchBookingsMap[listingID]

		allBookingsToCheck := append(dbBookings, batchBookings...)

		if err := checkTimeIntervals(&bookings[i], allBookingsToCheck); err != nil {
			return err
		}

		batchBookingsMap[listingID] = append(batchBookingsMap[listingID], bookings[i])

		bookings[i].HostID = dbListing.HostID
		bookings[i].TotalPrice = countTotalPrice(bookings[i].InDate, bookings[i].OutDate, dbListing.PricePerNight)
		bookings[i].IsPaid = false
	}

	return s.repo.CreateBookings(ctx, bookings)
}

func checkTimeIntervals(booking *model.Booking, bookings []model.Booking) error {
	for i := range bookings {

		if bookings[i].InDate.Before(booking.OutDate) && bookings[i].OutDate.After(booking.InDate) {
			zap.S().Errorf("bookings overlap")
			return errors.New("booking dates overlap with existing booking")
		}
	}

	return nil
}
