package service

import (
	"context"
	"time"

	"github.com/Rissochek/db-cw/internal/model"
)

func (s *Service) CreateBookingWithPayment(ctx context.Context, listingID, guestID int, inDate, outDate time.Time, paymentMethod string) (*model.CreateBookingWithPaymentResult, error) {
	return s.repo.CreateBookingWithPayment(ctx, listingID, guestID, inDate, outDate, paymentMethod)
}

func (s *Service) ConfirmPayment(ctx context.Context, paymentID int, transactionID *string) error {
	return s.repo.ConfirmPayment(ctx, paymentID, transactionID)
}

func (s *Service) CancelBookingWithRefund(ctx context.Context, bookingID int) error {
	return s.repo.CancelBookingWithRefund(ctx, bookingID)
}
