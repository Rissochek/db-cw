package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func (pg *Postgres) CreateBooking(ctx context.Context, booking *model.Booking) error {
	query := `INSERT INTO bookings (listing_id, host_id, guest_id, in_date, out_date, total_price, is_paid) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING booking_id`

	err := pg.conn.QueryRowxContext(ctx, query, booking.ListingID, booking.HostID, booking.GuestID,
		booking.InDate, booking.OutDate, booking.TotalPrice, booking.IsPaid).Scan(&booking.BookingID)
	if err != nil {
		zap.S().Errorf("failed to create booking: %v", err)
		return fmt.Errorf("failed to create booking")
	}

	return nil
}

func (pg *Postgres) CreateBookings(ctx context.Context, bookings []model.Booking) error {
	query := `INSERT INTO bookings (listing_id, host_id, guest_id, in_date, out_date, total_price, is_paid) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	zap.S().Infof("start adding %v bookings", len(bookings))
	tx, err := pg.conn.BeginTxx(ctx, nil)
	if err != nil {
		zap.S().Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("failed to create bookings")
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		zap.S().Errorf("failed to prepare statement: %v", err)
		return fmt.Errorf("failed to create bookings")
	}
	defer stmt.Close()

	for i := range bookings {
		_, err := stmt.ExecContext(ctx, bookings[i].ListingID, bookings[i].HostID, bookings[i].GuestID,
			bookings[i].InDate, bookings[i].OutDate, bookings[i].TotalPrice, bookings[i].IsPaid)
		if err != nil {
			zap.S().Errorf("failed to insert booking at index %d: %v", i, err)
			return fmt.Errorf("failed to create bookings")
		}
	}

	if err := tx.Commit(); err != nil {
		zap.S().Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to create bookings")
	}
	zap.S().Infof("added %v bookings", len(bookings))

	return nil
}

func (pg *Postgres) GetBookingByID(ctx context.Context, bookingID int) (*model.Booking, error) {
	var booking model.Booking

	query := `SELECT booking_id, listing_id, host_id, guest_id, in_date, out_date, total_price, is_paid FROM bookings WHERE booking_id = $1`

	err := pg.conn.GetContext(ctx, &booking, query, bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.S().Errorf("booking with id %d not found", bookingID)
			return nil, fmt.Errorf("booking not found")
		}
		zap.S().Errorf("failed to get booking: %v", err)
		return nil, fmt.Errorf("failed to get booking")
	}

	return &booking, nil
}

func (pg *Postgres) GetBookingsByID(ctx context.Context, bookingIDs []int) ([]model.Booking, error) {
	query, args, err := sqlx.In(`SELECT booking_id, listing_id, host_id, guest_id, in_date, out_date, total_price, is_paid FROM bookings WHERE booking_id IN (?)`, bookingIDs)
	if err != nil {
		zap.S().Errorf("failed to build query: %v", err)
		return nil, fmt.Errorf("failed to get bookings")
	}

	query = pg.conn.Rebind(query)
	var bookings []model.Booking

	err = pg.conn.SelectContext(ctx, &bookings, query, args...)
	if err != nil {
		zap.S().Errorf("failed to get bookings: %v", err)
		return nil, fmt.Errorf("failed to get bookings")
	}

	return bookings, nil
}

func (pg *Postgres) UpdateBooking(ctx context.Context, booking *model.Booking) error {
	query := `UPDATE bookings SET listing_id = $1, host_id = $2, guest_id = $3, in_date = $4, out_date = $5, total_price = $6, is_paid = $7 WHERE booking_id = $8`

	result, err := pg.conn.ExecContext(ctx, query, booking.ListingID, booking.HostID, booking.GuestID,
		booking.InDate, booking.OutDate, booking.TotalPrice, booking.IsPaid, booking.BookingID)
	if err != nil {
		zap.S().Errorf("failed to update booking: %v", err)
		return fmt.Errorf("failed to update booking")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to update booking")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("booking with id %d not found", booking.BookingID)
		return fmt.Errorf("booking not found")
	}

	return nil
}

func (pg *Postgres) UpdateBookings(ctx context.Context, bookings []model.Booking) error {
	query := `UPDATE bookings SET listing_id = $1, host_id = $2, guest_id = $3, in_date = $4, out_date = $5, total_price = $6, is_paid = $7 WHERE booking_id = $8`

	tx, err := pg.conn.BeginTxx(ctx, nil)
	if err != nil {
		zap.S().Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("failed to update bookings")
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		zap.S().Errorf("failed to prepare statement: %v", err)
		return fmt.Errorf("failed to update bookings")
	}
	defer stmt.Close()

	for i := range bookings {
		_, err := stmt.ExecContext(ctx, bookings[i].ListingID, bookings[i].HostID, bookings[i].GuestID,
			bookings[i].InDate, bookings[i].OutDate, bookings[i].TotalPrice, bookings[i].IsPaid, bookings[i].BookingID)
		if err != nil {
			zap.S().Errorf("failed to update booking at index %d: %v", i, err)
			return fmt.Errorf("failed to update bookings")
		}
	}

	if err := tx.Commit(); err != nil {
		zap.S().Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to update bookings")
	}

	return nil
}

func (pg *Postgres) DeleteBooking(ctx context.Context, bookingID int) error {
	query := `DELETE FROM bookings WHERE booking_id = $1`

	result, err := pg.conn.ExecContext(ctx, query, bookingID)
	if err != nil {
		zap.S().Errorf("failed to delete booking: %v", err)
		return fmt.Errorf("failed to delete booking")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to delete booking")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("booking with id %d not found", bookingID)
		return fmt.Errorf("booking not found")
	}

	return nil
}

func (pg *Postgres) DeleteBookings(ctx context.Context, bookingIDs []int) error {
	query, args, err := sqlx.In(`DELETE FROM bookings WHERE booking_id IN (?)`, bookingIDs)
	if err != nil {
		zap.S().Errorf("failed to build query: %v", err)
		return fmt.Errorf("failed to delete bookings")
	}

	query = pg.conn.Rebind(query)

	_, err = pg.conn.ExecContext(ctx, query, args...)
	if err != nil {
		zap.S().Errorf("failed to delete bookings: %v", err)
		return fmt.Errorf("failed to delete bookings")
	}

	return nil
}
