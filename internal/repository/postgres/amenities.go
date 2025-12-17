package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Rissochek/db-cw/internal/model"
	"go.uber.org/zap"
)

func (pg *Postgres) CreateAmenity(ctx context.Context, amenity *model.Amenity) error {
	query := `INSERT INTO amenities (name) VALUES ($1) RETURNING id`

	err := pg.conn.QueryRowxContext(ctx, query, amenity.Name).Scan(&amenity.ID)
	if err != nil {
		zap.S().Errorf("failed to create amenity: %v", err)
		return fmt.Errorf("failed to create amenity")
	}

	return nil
}

func (pg *Postgres) GetAmenityByID(ctx context.Context, id int) (*model.Amenity, error) {
	var amenity model.Amenity
	query := `SELECT id, name FROM amenities WHERE id = $1`
	err := pg.conn.GetContext(ctx, &amenity, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.S().Errorf("amenity with id %d not found", id)
			return nil, fmt.Errorf("amenity not found")
		}
		zap.S().Errorf("failed to get amenity: %v", err)
		return nil, fmt.Errorf("failed to get amenity")
	}
	return &amenity, nil
}

func (pg *Postgres) GetAllAmenities(ctx context.Context) ([]model.Amenity, error) {
	var amenities []model.Amenity
	query := `SELECT id, name FROM amenities ORDER BY id`
	err := pg.conn.SelectContext(ctx, &amenities, query)
	if err != nil {
		zap.S().Errorf("failed to get all amenities: %v", err)
		return nil, fmt.Errorf("failed to get all amenities")
	}
	return amenities, nil
}

func (pg *Postgres) UpdateAmenity(ctx context.Context, amenity *model.Amenity) error {
	query := `UPDATE amenities SET name = $1 WHERE id = $2`

	result, err := pg.conn.ExecContext(ctx, query, amenity.Name, amenity.ID)
	if err != nil {
		zap.S().Errorf("failed to update amenity: %v", err)
		return fmt.Errorf("failed to update amenity")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to update amenity")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("amenity with id %d not found", amenity.ID)
		return fmt.Errorf("amenity not found")
	}

	return nil
}

func (pg *Postgres) DeleteAmenity(ctx context.Context, id int) error {
	query := `DELETE FROM amenities WHERE id = $1`

	result, err := pg.conn.ExecContext(ctx, query, id)
	if err != nil {
		zap.S().Errorf("failed to delete amenity: %v", err)
		return fmt.Errorf("failed to delete amenity")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to delete amenity")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("amenity with id %d not found", id)
		return fmt.Errorf("amenity not found")
	}

	return nil
}

func (pg *Postgres) AddAmenityToListing(ctx context.Context, listingID int, amenityID int) error {
	query := `INSERT INTO listing_amenities (listing_id, amenity_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	_, err := pg.conn.ExecContext(ctx, query, listingID, amenityID)
	if err != nil {
		zap.S().Errorf("failed to add amenity %d to listing %d: %v", amenityID, listingID, err)
		return fmt.Errorf("failed to add amenity to listing")
	}

	return nil
}

func (pg *Postgres) RemoveAmenityFromListing(ctx context.Context, listingID int, amenityID int) error {
	query := `DELETE FROM listing_amenities WHERE listing_id = $1 AND amenity_id = $2`

	result, err := pg.conn.ExecContext(ctx, query, listingID, amenityID)
	if err != nil {
		zap.S().Errorf("failed to remove amenity %d from listing %d: %v", amenityID, listingID, err)
		return fmt.Errorf("failed to remove amenity from listing")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to remove amenity from listing")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("amenity %d not found for listing %d", amenityID, listingID)
		return fmt.Errorf("amenity not found for listing")
	}

	return nil
}

func (pg *Postgres) GetAmenitiesByListingID(ctx context.Context, listingID int) ([]model.Amenity, error) {
	var amenities []model.Amenity
	query := `SELECT a.id, a.name FROM amenities a
		INNER JOIN listing_amenities la ON a.id = la.amenity_id
		WHERE la.listing_id = $1
		ORDER BY a.id`
	err := pg.conn.SelectContext(ctx, &amenities, query, listingID)
	if err != nil {
		zap.S().Errorf("failed to get amenities for listing %d: %v", listingID, err)
		return nil, fmt.Errorf("failed to get amenities for listing")
	}
	return amenities, nil
}

func (pg *Postgres) CreateListingAmenities(ctx context.Context, listingAmenities []model.ListingAmenity) error {
	if len(listingAmenities) == 0 {
		return nil
	}

	query := `INSERT INTO listing_amenities (listing_id, amenity_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	zap.S().Infof("start adding %v listing amenities", len(listingAmenities))
	tx, err := pg.conn.BeginTxx(ctx, nil)
	if err != nil {
		zap.S().Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("failed to create listing amenities")
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		zap.S().Errorf("failed to prepare statement: %v", err)
		return fmt.Errorf("failed to create listing amenities")
	}
	defer stmt.Close()

	for i := range listingAmenities {
		_, err := stmt.ExecContext(ctx, listingAmenities[i].ListingID, listingAmenities[i].AmenityID)
		if err != nil {
			zap.S().Errorf("failed to insert listing amenity at index %d: %v", i, err)
			return fmt.Errorf("failed to create listing amenities")
		}
	}

	if err := tx.Commit(); err != nil {
		zap.S().Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to create listing amenities")
	}
	zap.S().Infof("added %v listing amenities", len(listingAmenities))

	return nil
}
