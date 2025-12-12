package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/jmoiron/sqlx"
)

func (pg *Postgres) CreateReview(ctx context.Context, review *model.Review) error {
	query := `INSERT INTO reviews (booking_id, user_id, text, score) VALUES ($1, $2, $3, $4) RETURNING id`

	err := pg.conn.QueryRowxContext(ctx, query, review.BookingID, review.UserID, review.Text, review.Score).Scan(&review.ID)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

func (pg *Postgres) CreateReviews(ctx context.Context, reviews []model.Review) error {
	query := `INSERT INTO reviews (booking_id, user_id, text, score) VALUES ($1, $2, $3, $4)`

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

	for i := range reviews {
		_, err := stmt.ExecContext(ctx, reviews[i].BookingID, reviews[i].UserID, reviews[i].Text, reviews[i].Score)
		if err != nil {
			return fmt.Errorf("failed to insert review at index %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (pg *Postgres) GetReviewByID(ctx context.Context, id int) (*model.Review, error) {
	var review model.Review

	query := `SELECT id, booking_id, user_id, text, score FROM reviews WHERE id = $1`

	err := pg.conn.GetContext(ctx, &review, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("review with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	return &review, nil
}

func (pg *Postgres) GetReviewsByID(ctx context.Context, ids []int) ([]model.Review, error) {

	query, args, err := sqlx.In(`SELECT id, booking_id, user_id, text, score FROM reviews WHERE id IN (?)`, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	query = pg.conn.Rebind(query)
	var reviews []model.Review

	err = pg.conn.SelectContext(ctx, &reviews, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	return reviews, nil
}

func (pg *Postgres) UpdateReview(ctx context.Context, review *model.Review) error {
	query := `UPDATE reviews
		SET booking_id = $1, user_id = $2, text = $3, score = $4
		WHERE id = $5`

	result, err := pg.conn.ExecContext(ctx, query, review.BookingID, review.UserID, review.Text, review.Score, review.ID)
	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("review with id %d not found", review.ID)
	}

	return nil
}

func (pg *Postgres) UpdateReviews(ctx context.Context, reviews []model.Review) error {
	query := `UPDATE reviews
		SET booking_id = $1, user_id = $2, text = $3, score = $4
		WHERE id = $5`

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

	for i := range reviews {
		_, err := stmt.ExecContext(ctx, reviews[i].BookingID, reviews[i].UserID, reviews[i].Text, reviews[i].Score, reviews[i].ID)
		if err != nil {
			return fmt.Errorf("failed to update review at index %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (pg *Postgres) DeleteReview(ctx context.Context, id int) error {
	query := `DELETE FROM reviews WHERE id = $1`

	result, err := pg.conn.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("review with id %d not found", id)
	}

	return nil
}

func (pg *Postgres) DeleteReviews(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`DELETE FROM reviews WHERE id IN (?)`, ids)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	query = pg.conn.Rebind(query)

	_, err = pg.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete reviews: %w", err)
	}

	return nil
}
