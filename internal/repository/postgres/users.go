package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/jmoiron/sqlx"
)

func (pg *Postgres) CreateUser(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (email, password, first_name, second_name) VALUES ($1, $2, $3, $4) RETURNING id`

	err := pg.conn.QueryRowxContext(ctx, query, user.Email, user.Password, user.FirstName, user.SecondName).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (pg *Postgres) CreateUsers(ctx context.Context, users []model.User) error {
	query := `INSERT INTO users (email, password, first_name, second_name) VALUES ($1, $2, $3, $4)`

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

	for i := range users {
		_, err := stmt.ExecContext(ctx, users[i].Email, users[i].Password, users[i].FirstName, users[i].SecondName)
		if err != nil {
			return fmt.Errorf("failed to insert user at index %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (pg *Postgres) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User

	query := `SELECT id, email, password, first_name, second_name FROM users WHERE id = $1`

	err := pg.conn.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (pg *Postgres) GetUsersByID(ctx context.Context, ids []int) ([]model.User, error) {
	query, args, err := sqlx.In(`SELECT id, email, password, first_name, second_name FROM users WHERE id IN (?)`, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	query = pg.conn.Rebind(query)
	var users []model.User

	err = pg.conn.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (pg *Postgres) UpdateUser(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET email = $1, password = $2, first_name = $3, second_name = $4 WHERE id = $5`

	result, err := pg.conn.ExecContext(ctx, query, user.Email, user.Password, user.FirstName, user.SecondName, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", user.ID)
	}

	return nil
}

func (pg *Postgres) UpdateUsers(ctx context.Context, users []model.User) error {
	query := `UPDATE users SET email = $1, password = $2, first_name = $3, second_name = $4 WHERE id = $5`

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

	for i := range users {
		_, err := stmt.ExecContext(ctx, users[i].Email, users[i].Password, users[i].FirstName, users[i].SecondName, users[i].ID)
		if err != nil {
			return fmt.Errorf("failed to update user at index %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (pg *Postgres) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := pg.conn.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

func (pg *Postgres) DeleteUsers(ctx context.Context, ids []int) error {
	query, args, err := sqlx.In(`DELETE FROM users WHERE id IN (?)`, ids)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	query = pg.conn.Rebind(query)

	_, err = pg.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete users: %w", err)
	}

	return nil
}
