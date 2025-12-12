package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func (pg *Postgres) CreateReview(ctx context.Context, review *model.Review) error {
	query := `INSERT INTO reviews (booking_id, user_id, text, score) VALUES ($1, $2, $3, $4) RETURNING id`

	err := pg.conn.QueryRowxContext(ctx, query, review.BookingID, review.UserID, review.Text, review.Score).Scan(&review.ID)
	if err != nil {
		zap.S().Errorf("failed to create review: %v", err)
		return fmt.Errorf("failed to create review")
	}

	return nil
}

func (pg *Postgres) CreateReviews(ctx context.Context, reviews []model.Review) error {
	query := `INSERT INTO reviews (booking_id, user_id, text, score) VALUES ($1, $2, $3, $4)`

	zap.S().Infof("start adding %v reviews", len(reviews))
	tx, err := pg.conn.BeginTxx(ctx, nil)
	if err != nil {
		zap.S().Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("failed to create reviews")
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		zap.S().Errorf("failed to prepare statement: %v", err)
		return fmt.Errorf("failed to create reviews")
	}
	defer stmt.Close()

	for i := range reviews {
		_, err := stmt.ExecContext(ctx, reviews[i].BookingID, reviews[i].UserID, reviews[i].Text, reviews[i].Score)
		if err != nil {
			zap.S().Errorf("failed to insert review at index %d: %v", i, err)
			return fmt.Errorf("failed to create reviews")
		}
	}

	if err := tx.Commit(); err != nil {
		zap.S().Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to create reviews")
	}
	zap.S().Infof("added %v reviews", len(reviews))

	return nil
}

func (pg *Postgres) GetReviewByID(ctx context.Context, id int) (*model.Review, error) {
	var review model.Review

	query := `SELECT id, booking_id, user_id, text, score FROM reviews WHERE id = $1`

	err := pg.conn.GetContext(ctx, &review, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.S().Errorf("review with id %d not found", id)
			return nil, fmt.Errorf("review not found")
		}
		zap.S().Errorf("failed to get review: %v", err)
		return nil, fmt.Errorf("failed to get review")
	}

	return &review, nil
}

func (pg *Postgres) GetReviewsByID(ctx context.Context, ids []int) ([]model.Review, error) {

	query, args, err := sqlx.In(`SELECT id, booking_id, user_id, text, score FROM reviews WHERE id IN (?)`, ids)
	if err != nil {
		zap.S().Errorf("failed to build query: %v", err)
		return nil, fmt.Errorf("failed to get reviews")
	}

	query = pg.conn.Rebind(query)
	var reviews []model.Review

	err = pg.conn.SelectContext(ctx, &reviews, query, args...)
	if err != nil {
		zap.S().Errorf("failed to get reviews: %v", err)
		return nil, fmt.Errorf("failed to get reviews")
	}

	return reviews, nil
}

func (pg *Postgres) UpdateReview(ctx context.Context, review *model.Review) error {
	query := `UPDATE reviews
		SET booking_id = $1, user_id = $2, text = $3, score = $4
		WHERE id = $5`

	result, err := pg.conn.ExecContext(ctx, query, review.BookingID, review.UserID, review.Text, review.Score, review.ID)
	if err != nil {
		zap.S().Errorf("failed to update review: %v", err)
		return fmt.Errorf("failed to update review")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to update review")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("review with id %d not found", review.ID)
		return fmt.Errorf("review not found")
	}

	return nil
}

func (pg *Postgres) UpdateReviews(ctx context.Context, reviews []model.Review) error {
	query := `UPDATE reviews
		SET booking_id = $1, user_id = $2, text = $3, score = $4
		WHERE id = $5`

	tx, err := pg.conn.BeginTxx(ctx, nil)
	if err != nil {
		zap.S().Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("failed to update reviews")
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		zap.S().Errorf("failed to prepare statement: %v", err)
		return fmt.Errorf("failed to update reviews")
	}
	defer stmt.Close()

	for i := range reviews {
		_, err := stmt.ExecContext(ctx, reviews[i].BookingID, reviews[i].UserID, reviews[i].Text, reviews[i].Score, reviews[i].ID)
		if err != nil {
			zap.S().Errorf("failed to update review at index %d: %v", i, err)
			return fmt.Errorf("failed to update reviews")
		}
	}

	if err := tx.Commit(); err != nil {
		zap.S().Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to update reviews")
	}

	return nil
}

func (pg *Postgres) DeleteReview(ctx context.Context, id int) error {
	query := `DELETE FROM reviews WHERE id = $1`

	result, err := pg.conn.ExecContext(ctx, query, id)
	if err != nil {
		zap.S().Errorf("failed to delete review: %v", err)
		return fmt.Errorf("failed to delete review")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to delete review")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("review with id %d not found", id)
		return fmt.Errorf("review not found")
	}

	return nil
}

func (pg *Postgres) DeleteReviews(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`DELETE FROM reviews WHERE id IN (?)`, ids)
	if err != nil {
		zap.S().Errorf("failed to build query: %v", err)
		return fmt.Errorf("failed to delete reviews")
	}

	query = pg.conn.Rebind(query)

	_, err = pg.conn.ExecContext(ctx, query, args...)
	if err != nil {
		zap.S().Errorf("failed to delete reviews: %v", err)
		return fmt.Errorf("failed to delete reviews")
	}

	return nil
}
