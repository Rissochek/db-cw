package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Rissochek/db-cw/internal/model"
	"go.uber.org/zap"
)

func (pg *Postgres) CreateFavorite(ctx context.Context, favorite *model.Favorite) error {
	query := `INSERT INTO favorites (user_id, listing_id) VALUES ($1, $2) RETURNING id`

	err := pg.conn.QueryRowxContext(ctx, query, favorite.UserID, favorite.ListingID).Scan(&favorite.ID)
	if err != nil {
		zap.S().Errorf("failed to create favorite: %v", err)
		return fmt.Errorf("failed to create favorite")
	}

	return nil
}

func (pg *Postgres) GetFavoriteByID(ctx context.Context, id int) (*model.Favorite, error) {
	var favorite model.Favorite
	query := `SELECT id, user_id, listing_id FROM favorites WHERE id = $1`
	err := pg.conn.GetContext(ctx, &favorite, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.S().Errorf("favorite with id %d not found", id)
			return nil, fmt.Errorf("favorite not found")
		}
		zap.S().Errorf("failed to get favorite: %v", err)
		return nil, fmt.Errorf("failed to get favorite")
	}
	return &favorite, nil
}

func (pg *Postgres) GetFavoritesByUserID(ctx context.Context, userID int) ([]model.Favorite, error) {
	var favorites []model.Favorite
	query := `SELECT id, user_id, listing_id FROM favorites WHERE user_id = $1 ORDER BY id`
	err := pg.conn.SelectContext(ctx, &favorites, query, userID)
	if err != nil {
		zap.S().Errorf("failed to get favorites for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to get favorites for user")
	}
	return favorites, nil
}

func (pg *Postgres) GetFavoriteByUserAndListing(ctx context.Context, userID int, listingID int) (*model.Favorite, error) {
	var favorite model.Favorite
	query := `SELECT id, user_id, listing_id FROM favorites WHERE user_id = $1 AND listing_id = $2`
	err := pg.conn.GetContext(ctx, &favorite, query, userID, listingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		zap.S().Errorf("failed to get favorite: %v", err)
		return nil, fmt.Errorf("failed to get favorite")
	}
	return &favorite, nil
}

func (pg *Postgres) DeleteFavorite(ctx context.Context, id int) error {
	query := `DELETE FROM favorites WHERE id = $1`

	result, err := pg.conn.ExecContext(ctx, query, id)
	if err != nil {
		zap.S().Errorf("failed to delete favorite: %v", err)
		return fmt.Errorf("failed to delete favorite")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to delete favorite")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("favorite with id %d not found", id)
		return fmt.Errorf("favorite not found")
	}

	return nil
}

func (pg *Postgres) DeleteFavoriteByUserAndListing(ctx context.Context, userID int, listingID int) error {
	query := `DELETE FROM favorites WHERE user_id = $1 AND listing_id = $2`

	result, err := pg.conn.ExecContext(ctx, query, userID, listingID)
	if err != nil {
		zap.S().Errorf("failed to delete favorite: %v", err)
		return fmt.Errorf("failed to delete favorite")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to delete favorite")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("favorite for user %d and listing %d not found", userID, listingID)
		return fmt.Errorf("favorite not found")
	}

	return nil
}

func (pg *Postgres) CreateFavorites(ctx context.Context, favorites []model.Favorite) error {
	if len(favorites) == 0 {
		return nil
	}

	query := `INSERT INTO favorites (user_id, listing_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	zap.S().Infof("start adding %v favorites", len(favorites))
	tx, err := pg.conn.BeginTxx(ctx, nil)
	if err != nil {
		zap.S().Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("failed to create favorites")
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		zap.S().Errorf("failed to prepare statement: %v", err)
		return fmt.Errorf("failed to create favorites")
	}
	defer stmt.Close()

	for i := range favorites {
		_, err := stmt.ExecContext(ctx, favorites[i].UserID, favorites[i].ListingID)
		if err != nil {
			zap.S().Errorf("failed to insert favorite at index %d: %v", i, err)
			return fmt.Errorf("failed to create favorites")
		}
	}

	if err := tx.Commit(); err != nil {
		zap.S().Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to create favorites")
	}
	zap.S().Infof("added %v favorites", len(favorites))

	return nil
}
