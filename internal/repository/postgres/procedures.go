package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Rissochek/db-cw/internal/model"
	"go.uber.org/zap"
)

func (pg *Postgres) CreateBookingWithPayment(ctx context.Context, listingID, guestID int, inDate, outDate time.Time, paymentMethod string) (*model.CreateBookingWithPaymentResult, error) {
	var result model.CreateBookingWithPaymentResult
	var hostID int

	getHostQuery := `SELECT host_id FROM listings WHERE id = $1`
	err := pg.conn.GetContext(ctx, &hostID, getHostQuery, listingID)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.S().Errorf("listing with id %d not found", listingID)
			return nil, fmt.Errorf("listing not found")
		}
		zap.S().Errorf("failed to get host_id for listing %d: %v", listingID, err)
		return nil, fmt.Errorf("failed to get listing host")
	}

	query := `SELECT p_booking_id, p_payment_id FROM create_booking_with_payment($1, $2, $3, $4, $5, $6)`

	tx, err := pg.conn.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	})
	if err != nil {
		zap.S().Errorf("failed to begin transaction: %v", err)
		return nil, fmt.Errorf("failed to create booking with payment: %w", err)
	}
	defer tx.Rollback()

	err = tx.QueryRowxContext(ctx, query,
		listingID, hostID, guestID, inDate, outDate, paymentMethod).Scan(&result.BookingID, &result.PaymentID)
	if err != nil {
		zap.S().Errorf("failed to create booking with payment: %v", err)
		return nil, fmt.Errorf("failed to create booking with payment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		zap.S().Errorf("failed to commit transaction: %v", err)
		return nil, fmt.Errorf("failed to create booking with payment: %w", err)
	}

	return &result, nil
}

func (pg *Postgres) ConfirmPayment(ctx context.Context, paymentID int, transactionID *string) error {
	query := `CALL confirm_payment($1, $2)`

	tx, err := pg.conn.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	})
	if err != nil {
		zap.S().Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("failed to confirm payment: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, paymentID, transactionID)
	if err != nil {
		zap.S().Errorf("failed to confirm payment %d: %v", paymentID, err)
		return fmt.Errorf("failed to confirm payment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		zap.S().Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to confirm payment: %w", err)
	}

	return nil
}

func (pg *Postgres) CancelBookingWithRefund(ctx context.Context, bookingID int) error {
	query := `CALL cancel_booking_with_refund($1)`

	tx, err := pg.conn.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	})
	if err != nil {
		zap.S().Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("failed to cancel booking with refund: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, bookingID)
	if err != nil {
		zap.S().Errorf("failed to cancel booking with refund %d: %v", bookingID, err)
		return fmt.Errorf("failed to cancel booking with refund: %w", err)
	}

	if err := tx.Commit(); err != nil {
		zap.S().Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to cancel booking with refund: %w", err)
	}

	return nil
}
