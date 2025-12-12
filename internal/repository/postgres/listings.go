package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/jmoiron/sqlx"
)

func (pg *Postgres) CreateListing(ctx context.Context, listing *model.Listing) error {
	query := `INSERT INTO listings (host_id, address, price_per_night, is_available, rooms_number, beds_number) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := pg.conn.QueryRowxContext(ctx, query, listing.HostID, listing.Address, listing.PricePerNight,
		listing.IsAvailable, listing.RoomsNumber, listing.BedsNumber).Scan(&listing.ID)
	if err != nil {
		return fmt.Errorf("failed to create listing: %w", err)
	}

	return nil
}

func (pg *Postgres) CreateListings(ctx context.Context, listings []model.Listing) error {
	query := `INSERT INTO listings (host_id, address, price_per_night, is_available, rooms_number, beds_number) 
		VALUES ($1, $2, $3, $4, $5, $6)`

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

	for i := range listings {
		_, err := stmt.ExecContext(ctx, listings[i].HostID, listings[i].Address, listings[i].PricePerNight,
			listings[i].IsAvailable, listings[i].RoomsNumber, listings[i].BedsNumber)
		if err != nil {
			return fmt.Errorf("failed to insert listing at index %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (pg *Postgres) GetListingByID(ctx context.Context, id int) (*model.Listing, error) {
	var listing model.Listing

	query := `SELECT id, host_id, address, price_per_night, is_available, rooms_number, beds_number 
		FROM listings WHERE id = $1`

	err := pg.conn.GetContext(ctx, &listing, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("listing with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	return &listing, nil
}

func (pg *Postgres) GetListingsByID(ctx context.Context, ids []int) ([]model.Listing, error) {
	query, args, err := sqlx.In(`SELECT id, host_id, address, price_per_night, is_available, rooms_number, beds_number 
		FROM listings WHERE id IN (?)`, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	query = pg.conn.Rebind(query)
	var listings []model.Listing

	err = pg.conn.SelectContext(ctx, &listings, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get listings: %w", err)
	}

	return listings, nil
}

func (pg *Postgres) UpdateListing(ctx context.Context, listing *model.Listing) error {
	query := `UPDATE listings
		SET host_id = $1, address = $2, price_per_night = $3, is_available = $4, rooms_number = $5, beds_number = $6
		WHERE id = $7`

	result, err := pg.conn.ExecContext(ctx, query, listing.HostID, listing.Address, listing.PricePerNight,
		listing.IsAvailable, listing.RoomsNumber, listing.BedsNumber, listing.ID)
	if err != nil {
		return fmt.Errorf("failed to update listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("listing with id %d not found", listing.ID)
	}

	return nil
}

func (pg *Postgres) UpdateListings(ctx context.Context, listings []model.Listing) error {
	query := `UPDATE listings SET host_id = $1, address = $2, price_per_night = $3, is_available = $4, rooms_number = $5, beds_number = $6 WHERE id = $7`

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

	for i := range listings {
		_, err := stmt.ExecContext(ctx, listings[i].HostID, listings[i].Address, listings[i].PricePerNight,
			listings[i].IsAvailable, listings[i].RoomsNumber, listings[i].BedsNumber, listings[i].ID)
		if err != nil {
			return fmt.Errorf("failed to update listing at index %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (pg *Postgres) DeleteListing(ctx context.Context, id int) error {
	query := `DELETE FROM listings WHERE id = $1`

	result, err := pg.conn.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("listing with id %d not found", id)
	}

	return nil
}

func (pg *Postgres) DeleteListings(ctx context.Context, ids []int) error {
	query, args, err := sqlx.In(`DELETE FROM listings WHERE id IN (?)`, ids)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	query = pg.conn.Rebind(query)

	_, err = pg.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete listings: %w", err)
	}

	return nil
}
