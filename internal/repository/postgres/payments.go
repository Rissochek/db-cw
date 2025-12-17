package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Rissochek/db-cw/internal/model"
	"go.uber.org/zap"
)

func (pg *Postgres) CreatePayment(ctx context.Context, payment *model.Payment) error {
	query := `INSERT INTO payments (booking_id, amount, payment_method, payment_status, transaction_id, paid_at) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING payment_id`

	err := pg.conn.QueryRowxContext(ctx, query, payment.BookingID, payment.Amount, payment.PaymentMethod,
		payment.PaymentStatus, payment.TransactionID, payment.PaidAt).Scan(&payment.PaymentID)
	if err != nil {
		zap.S().Errorf("failed to create payment: %v", err)
		return fmt.Errorf("failed to create payment")
	}

	return nil
}

func (pg *Postgres) GetPaymentByID(ctx context.Context, paymentID int) (*model.Payment, error) {
	var payment model.Payment
	query := `SELECT payment_id, booking_id, amount, payment_method, payment_status, transaction_id, paid_at 
		FROM payments WHERE payment_id = $1`
	err := pg.conn.GetContext(ctx, &payment, query, paymentID)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.S().Errorf("payment with id %d not found", paymentID)
			return nil, fmt.Errorf("payment not found")
		}
		zap.S().Errorf("failed to get payment: %v", err)
		return nil, fmt.Errorf("failed to get payment")
	}
	return &payment, nil
}

func (pg *Postgres) GetPaymentsByBookingID(ctx context.Context, bookingID int) ([]model.Payment, error) {
	var payments []model.Payment
	query := `SELECT payment_id, booking_id, amount, payment_method, payment_status, transaction_id, paid_at 
		FROM payments WHERE booking_id = $1 ORDER BY payment_id`
	err := pg.conn.SelectContext(ctx, &payments, query, bookingID)
	if err != nil {
		zap.S().Errorf("failed to get payments for booking %d: %v", bookingID, err)
		return nil, fmt.Errorf("failed to get payments for booking")
	}
	return payments, nil
}

func (pg *Postgres) UpdatePayment(ctx context.Context, payment *model.Payment) error {
	query := `UPDATE payments 
		SET booking_id = $1, amount = $2, payment_method = $3, payment_status = $4, transaction_id = $5, paid_at = $6 
		WHERE payment_id = $7`

	result, err := pg.conn.ExecContext(ctx, query, payment.BookingID, payment.Amount, payment.PaymentMethod,
		payment.PaymentStatus, payment.TransactionID, payment.PaidAt, payment.PaymentID)
	if err != nil {
		zap.S().Errorf("failed to update payment: %v", err)
		return fmt.Errorf("failed to update payment")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to update payment")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("payment with id %d not found", payment.PaymentID)
		return fmt.Errorf("payment not found")
	}

	return nil
}

func (pg *Postgres) DeletePayment(ctx context.Context, paymentID int) error {
	query := `DELETE FROM payments WHERE payment_id = $1`

	result, err := pg.conn.ExecContext(ctx, query, paymentID)
	if err != nil {
		zap.S().Errorf("failed to delete payment: %v", err)
		return fmt.Errorf("failed to delete payment")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to delete payment")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("payment with id %d not found", paymentID)
		return fmt.Errorf("payment not found")
	}

	return nil
}

func (pg *Postgres) CreatePayments(ctx context.Context, payments []model.Payment) error {
	if len(payments) == 0 {
		return nil
	}

	query := `INSERT INTO payments (booking_id, amount, payment_method, payment_status, transaction_id, paid_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	zap.S().Infof("start adding %v payments", len(payments))
	tx, err := pg.conn.BeginTxx(ctx, nil)
	if err != nil {
		zap.S().Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("failed to create payments")
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		zap.S().Errorf("failed to prepare statement: %v", err)
		return fmt.Errorf("failed to create payments")
	}
	defer stmt.Close()

	for i := range payments {
		_, err := stmt.ExecContext(ctx, payments[i].BookingID, payments[i].Amount, payments[i].PaymentMethod,
			payments[i].PaymentStatus, payments[i].TransactionID, payments[i].PaidAt)
		if err != nil {
			zap.S().Errorf("failed to insert payment at index %d: %v", i, err)
			return fmt.Errorf("failed to create payments")
		}
	}

	if err := tx.Commit(); err != nil {
		zap.S().Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to create payments")
	}
	zap.S().Infof("added %v payments", len(payments))

	return nil
}
