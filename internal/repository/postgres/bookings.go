package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/jmoiron/sqlx"
)

func (pg *Postgres) CreateBooking(ctx context.Context, booking *model.Booking) error {
	query := `INSERT INTO bookings (listing_id, host_id, guest_id, in_date, out_date, total_price, is_paid) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING booking_id`

	err := pg.conn.QueryRowxContext(ctx, query, booking.ListingID, booking.HostID, booking.GuestID,
		booking.InDate, booking.OutDate, booking.TotalPrice, booking.IsPaid).Scan(&booking.BookingID)
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	return nil
}

func (pg *Postgres) CreateBookings(ctx context.Context, bookings []model.Booking) error {
	query := `INSERT INTO bookings (listing_id, host_id, guest_id, in_date, out_date, total_price, is_paid) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	tx, err := pg.conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i := range bookings {
		_, err := stmt.ExecContext(ctx, bookings[i].ListingID, bookings[i].HostID, bookings[i].GuestID,
			bookings[i].InDate, bookings[i].OutDate, bookings[i].TotalPrice, bookings[i].IsPaid)
		if err != nil {
			return fmt.Errorf("failed to insert booking at index %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (pg *Postgres) GetBookingByID(ctx context.Context, bookingID int) (*model.Booking, error) {
	var booking model.Booking

	query := `SELECT booking_id, listing_id, host_id, guest_id, in_date, out_date, total_price, is_paid FROM bookings WHERE booking_id = $1`

	err := pg.conn.GetContext(ctx, &booking, query, bookingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking with id %d not found", bookingID)
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return &booking, nil
}

func (pg *Postgres) GetBookingsByID(ctx context.Context, bookingIDs []int) ([]model.Booking, error) {
	query, args, err := sqlx.In(`SELECT booking_id, listing_id, host_id, guest_id, in_date, out_date, total_price, is_paid FROM bookings WHERE booking_id IN (?)`, bookingIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	query = pg.conn.Rebind(query)
	var bookings []model.Booking

	err = pg.conn.SelectContext(ctx, &bookings, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	return bookings, nil
}

func (pg *Postgres) UpdateBooking(ctx context.Context, booking *model.Booking) error {
	query := `UPDATE bookings SET listing_id = $1, host_id = $2, guest_id = $3, in_date = $4, out_date = $5, total_price = $6, is_paid = $7 WHERE booking_id = $8`

	result, err := pg.conn.ExecContext(ctx, query, booking.ListingID, booking.HostID, booking.GuestID,
		booking.InDate, booking.OutDate, booking.TotalPrice, booking.IsPaid, booking.BookingID)
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking with id %d not found", booking.BookingID)
	}

	return nil
}

func (pg *Postgres) UpdateBookings(ctx context.Context, bookings []model.Booking) error {
	query := `UPDATE bookings SET listing_id = $1, host_id = $2, guest_id = $3, in_date = $4, out_date = $5, total_price = $6, is_paid = $7 WHERE booking_id = $8`

	tx, err := pg.conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i := range bookings {
		_, err := stmt.ExecContext(ctx, bookings[i].ListingID, bookings[i].HostID, bookings[i].GuestID,
			bookings[i].InDate, bookings[i].OutDate, bookings[i].TotalPrice, bookings[i].IsPaid, bookings[i].BookingID)
		if err != nil {
			return fmt.Errorf("failed to update booking at index %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (pg *Postgres) DeleteBooking(ctx context.Context, bookingID int) error {
	query := `DELETE FROM bookings WHERE booking_id = $1`

	result, err := pg.conn.ExecContext(ctx, query, bookingID)
	if err != nil {
		return fmt.Errorf("failed to delete booking: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking with id %d not found", bookingID)
	}

	return nil
}

func (pg *Postgres) DeleteBookings(ctx context.Context, bookingIDs []int) error {
	query, args, err := sqlx.In(`DELETE FROM bookings WHERE booking_id IN (?)`, bookingIDs)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	query = pg.conn.Rebind(query)

	_, err = pg.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete bookings: %w", err)
	}

	return nil
}
