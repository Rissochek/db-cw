package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Rissochek/db-cw/internal/faking"
	"github.com/Rissochek/db-cw/internal/model"
)

func (s *Service) CreatePayment(ctx context.Context, payment *model.Payment) error {
	booking, err := s.repo.GetBookingByID(ctx, payment.BookingID)
	if err != nil {
		return err
	}

	if payment.PaymentStatus == "completed" {
		payment.Amount = booking.TotalPrice

		if payment.PaidAt == nil {
			now := time.Now()
			payment.PaidAt = &now
		}

		if payment.TransactionID == nil {
			transactionID := faking.GenerateTransactionID(payment.BookingID)
			payment.TransactionID = &transactionID
		}
	} else if payment.PaymentStatus == "failed" {
		payment.Amount = 0
		payment.PaidAt = nil
		if payment.TransactionID != nil {
			payment.TransactionID = nil
		}
	} else {
		if payment.Amount <= 0 {
			return fmt.Errorf("payment amount must be greater than zero")
		}

		if payment.Amount > booking.TotalPrice {
			return fmt.Errorf("payment amount cannot exceed booking total price")
		}
	}

	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return err
	}

	return s.updateBookingIsPaidStatus(ctx, payment.BookingID)
}

func (s *Service) GetPaymentByID(ctx context.Context, paymentID int) (*model.Payment, error) {
	return s.repo.GetPaymentByID(ctx, paymentID)
}

func (s *Service) GetPaymentsByBookingID(ctx context.Context, bookingID int) ([]model.Payment, error) {
	return s.repo.GetPaymentsByBookingID(ctx, bookingID)
}

func (s *Service) UpdatePayment(ctx context.Context, payment *model.Payment) error {
	dbPayment, err := s.repo.GetPaymentByID(ctx, payment.PaymentID)
	if err != nil {
		return err
	}

	payment.BookingID = dbPayment.BookingID

	booking, err := s.repo.GetBookingByID(ctx, payment.BookingID)
	if err != nil {
		return err
	}

	if payment.PaymentStatus == "completed" {
		payment.Amount = booking.TotalPrice

		if payment.PaidAt == nil {
			now := time.Now()
			payment.PaidAt = &now
		}

		if payment.TransactionID == nil {
			transactionID := faking.GenerateTransactionID(payment.BookingID)
			payment.TransactionID = &transactionID
		}
	} else if payment.PaymentStatus == "failed" {
		payment.Amount = 0
		payment.PaidAt = nil
		payment.TransactionID = nil
	} else {
		payment.Amount = dbPayment.Amount
		if payment.PaidAt == nil {
			payment.PaidAt = dbPayment.PaidAt
		}
		if payment.TransactionID == nil {
			payment.TransactionID = dbPayment.TransactionID
		}
	}

	if err := s.repo.UpdatePayment(ctx, payment); err != nil {
		return err
	}

	return s.updateBookingIsPaidStatus(ctx, payment.BookingID)
}

func (s *Service) DeletePayment(ctx context.Context, paymentID int) error {
	payment, err := s.repo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return err
	}

	bookingID := payment.BookingID

	if err := s.repo.DeletePayment(ctx, paymentID); err != nil {
		return err
	}

	return s.updateBookingIsPaidStatus(ctx, bookingID)
}

func (s *Service) CreatePayments(ctx context.Context, payments []model.Payment) error {
	if err := s.repo.CreatePayments(ctx, payments); err != nil {
		return err
	}

	bookingIDs := make(map[int]bool)
	for i := range payments {
		bookingIDs[payments[i].BookingID] = true
	}

	for bookingID := range bookingIDs {
		if err := s.updateBookingIsPaidStatus(ctx, bookingID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) updateBookingIsPaidStatus(ctx context.Context, bookingID int) error {
	payments, err := s.repo.GetPaymentsByBookingID(ctx, bookingID)
	if err != nil {
		return err
	}

	hasCompleted := false
	for i := range payments {
		if payments[i].PaymentStatus == "completed" {
			hasCompleted = true
			break
		}
	}

	return s.repo.UpdateBookingIsPaid(ctx, bookingID, hasCompleted)
}

func (s *Service) UpdateBookingIsPaid(ctx context.Context, bookingID int, isPaid bool) error {
	return s.repo.UpdateBookingIsPaid(ctx, bookingID, isPaid)
}
