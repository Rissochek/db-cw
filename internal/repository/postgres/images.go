package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Rissochek/db-cw/internal/model"
	"go.uber.org/zap"
)

func (pg *Postgres) CreateImage(ctx context.Context, image *model.Image) error {
	query := `INSERT INTO images (listing_id, image_url, is_primary, order_index, uploaded_at) 
		VALUES ($1, $2, $3, $4, $5) RETURNING image_id`

	err := pg.conn.QueryRowxContext(ctx, query, image.ListingID, image.ImageURL, image.IsPrimary,
		image.OrderIndex, image.UploadedAt).Scan(&image.ImageID)
	if err != nil {
		zap.S().Errorf("failed to create image: %v", err)
		return fmt.Errorf("failed to create image")
	}

	return nil
}

func (pg *Postgres) GetImageByID(ctx context.Context, imageID int) (*model.Image, error) {
	var image model.Image
	query := `SELECT image_id, listing_id, image_url, is_primary, order_index, uploaded_at 
		FROM images WHERE image_id = $1`
	err := pg.conn.GetContext(ctx, &image, query, imageID)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.S().Errorf("image with id %d not found", imageID)
			return nil, fmt.Errorf("image not found")
		}
		zap.S().Errorf("failed to get image: %v", err)
		return nil, fmt.Errorf("failed to get image")
	}
	return &image, nil
}

func (pg *Postgres) GetImagesByListingID(ctx context.Context, listingID int) ([]model.Image, error) {
	var images []model.Image
	query := `SELECT image_id, listing_id, image_url, is_primary, order_index, uploaded_at 
		FROM images WHERE listing_id = $1 ORDER BY order_index, image_id`
	err := pg.conn.SelectContext(ctx, &images, query, listingID)
	if err != nil {
		zap.S().Errorf("failed to get images for listing %d: %v", listingID, err)
		return nil, fmt.Errorf("failed to get images for listing")
	}
	return images, nil
}

func (pg *Postgres) UpdateImage(ctx context.Context, image *model.Image) error {
	query := `UPDATE images 
		SET listing_id = $1, image_url = $2, is_primary = $3, order_index = $4, uploaded_at = $5 
		WHERE image_id = $6`

	result, err := pg.conn.ExecContext(ctx, query, image.ListingID, image.ImageURL, image.IsPrimary,
		image.OrderIndex, image.UploadedAt, image.ImageID)
	if err != nil {
		zap.S().Errorf("failed to update image: %v", err)
		return fmt.Errorf("failed to update image")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to update image")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("image with id %d not found", image.ImageID)
		return fmt.Errorf("image not found")
	}

	return nil
}

func (pg *Postgres) DeleteImage(ctx context.Context, imageID int) error {
	query := `DELETE FROM images WHERE image_id = $1`

	result, err := pg.conn.ExecContext(ctx, query, imageID)
	if err != nil {
		zap.S().Errorf("failed to delete image: %v", err)
		return fmt.Errorf("failed to delete image")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to delete image")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("image with id %d not found", imageID)
		return fmt.Errorf("image not found")
	}

	return nil
}

func (pg *Postgres) CreateImages(ctx context.Context, images []model.Image) error {
	if len(images) == 0 {
		return nil
	}

	query := `INSERT INTO images (listing_id, image_url, is_primary, order_index, uploaded_at) 
		VALUES ($1, $2, $3, $4, $5)`

	zap.S().Infof("start adding %v images", len(images))
	tx, err := pg.conn.BeginTxx(ctx, nil)
	if err != nil {
		zap.S().Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("failed to create images")
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		zap.S().Errorf("failed to prepare statement: %v", err)
		return fmt.Errorf("failed to create images")
	}
	defer stmt.Close()

	for i := range images {
		_, err := stmt.ExecContext(ctx, images[i].ListingID, images[i].ImageURL, images[i].IsPrimary,
			images[i].OrderIndex, images[i].UploadedAt)
		if err != nil {
			zap.S().Errorf("failed to insert image at index %d: %v", i, err)
			return fmt.Errorf("failed to create images")
		}
	}

	if err := tx.Commit(); err != nil {
		zap.S().Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to create images")
	}
	zap.S().Infof("added %v images", len(images))

	return nil
}
