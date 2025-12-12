package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Rissochek/db-cw/internal/model"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func (pg *Postgres) CreateUser(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (email, password, first_name, second_name) VALUES ($1, $2, $3, $4) RETURNING id`

	err := pg.conn.QueryRowxContext(ctx, query, user.Email, user.Password, user.FirstName, user.SecondName).Scan(&user.ID)
	if err != nil {
		zap.S().Errorf("failed to create user: %v", err)
		return fmt.Errorf("failed to create user")
	}

	return nil
}

func (pg *Postgres) CreateUsers(ctx context.Context, users []model.User) error {
	query := `INSERT INTO users (email, password, first_name, second_name) VALUES ($1, $2, $3, $4)`

	zap.S().Infof("start adding %v users", len(users))
	tx, err := pg.conn.BeginTxx(ctx, nil)
	if err != nil {
		zap.S().Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("failed to create users")
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		zap.S().Errorf("failed to prepare statement: %v", err)
		return fmt.Errorf("failed to create users")
	}
	defer stmt.Close()

	for i := range users {
		_, err := stmt.ExecContext(ctx, users[i].Email, users[i].Password, users[i].FirstName, users[i].SecondName)
		if err != nil {
			zap.S().Errorf("failed to insert user at index %d: %v", i, err)
			return fmt.Errorf("failed to create users")
		}
	}

	if err := tx.Commit(); err != nil {
		zap.S().Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to create users")
	}
	zap.S().Infof("added %v users", len(users))

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
		zap.S().Errorf("failed to get user: %v", err)
		return nil, fmt.Errorf("failed to get user")
	}

	return &user, nil
}

func (pg *Postgres) GetUsersByID(ctx context.Context, ids []int) ([]model.User, error) {
	query, args, err := sqlx.In(`SELECT id, email, password, first_name, second_name FROM users WHERE id IN (?)`, ids)
	if err != nil {
		zap.S().Errorf("failed to build query: %v", err)
		return nil, fmt.Errorf("failed to get users")
	}

	query = pg.conn.Rebind(query)
	var users []model.User

	err = pg.conn.SelectContext(ctx, &users, query, args...)
	if err != nil {
		zap.S().Errorf("failed to get users: %v", err)
		return nil, fmt.Errorf("failed to get users")
	}

	return users, nil
}

func (pg *Postgres) UpdateUser(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET email = $1, password = $2, first_name = $3, second_name = $4 WHERE id = $5`

	result, err := pg.conn.ExecContext(ctx, query, user.Email, user.Password, user.FirstName, user.SecondName, user.ID)
	if err != nil {
		zap.S().Errorf("failed to update user: %v", err)
		return fmt.Errorf("failed to update user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to update user")
	}

	if rowsAffected == 0 {
		zap.S().Errorf("user with id %d not found", user.ID)
		return fmt.Errorf("user not found")
	}

	return nil
}

func (pg *Postgres) UpdateUsers(ctx context.Context, users []model.User) error {
	query := `UPDATE users SET email = $1, password = $2, first_name = $3, second_name = $4 WHERE id = $5`

	tx, err := pg.conn.BeginTxx(ctx, nil)
	if err != nil {
		zap.S().Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("failed to update users")
	}
	defer tx.Rollback()

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		zap.S().Errorf("failed to prepare statement: %v", err)
		return fmt.Errorf("failed to update users")
	}
	defer stmt.Close()

	for i := range users {
		_, err := stmt.ExecContext(ctx, users[i].Email, users[i].Password, users[i].FirstName, users[i].SecondName, users[i].ID)
		if err != nil {
			zap.S().Errorf("failed to update user at index %d: %v", i, err)
			return fmt.Errorf("failed to update users")
		}
	}

	if err := tx.Commit(); err != nil {
		zap.S().Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to update users")
	}

	return nil
}

func (pg *Postgres) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := pg.conn.ExecContext(ctx, query, id)
	if err != nil {
		zap.S().Errorf("failed to delete user: %v", err)
		return fmt.Errorf("failed to delete user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.S().Errorf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to delete user")
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

func (pg *Postgres) DeleteUsers(ctx context.Context, ids []int) error {
	query, args, err := sqlx.In(`DELETE FROM users WHERE id IN (?)`, ids)
	if err != nil {
		zap.S().Errorf("failed to build query: %v", err)
		return fmt.Errorf("failed to delete users")
	}

	query = pg.conn.Rebind(query)

	_, err = pg.conn.ExecContext(ctx, query, args...)
	if err != nil {
		zap.S().Errorf("failed to delete users: %v", err)
		return fmt.Errorf("failed to delete users")
	}

	return nil
}
