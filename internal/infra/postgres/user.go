package db

import (
	"context"
	"database/sql"
	"iziumov/tv-v-x/internal/domain"
	"iziumov/tv-v-x/internal/dto"
)

type UserRepo struct {
	*DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) GetByTgId(ctx context.Context, telegram_id int64) (*domain.User, error) {
	query := `
		SELECT tg_id, username, first_name, last_name, is_banned
		FROM users
		WHERE tg_id = $1
	`

	var u domain.User
	err := r.QueryRowContext(ctx, query, telegram_id).Scan(
		&u.ID,
		&u.Username,
		&u.FirstName,
		&u.LastName,
		&u.IsBanned,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, user dto.CreateUser) error {
	query := `
		INSERT INTO users (tg_id, username, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(tg_id) DO NOTHING
	`

	_, err := r.DB.ExecContext(ctx, query, user.TGID, user.Username, user.FirstName, user.LastName)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Ban(ctx context.Context, telegram_id int64) error {
	query := `
		UPDATE users
		SET is_banned = TRUE
		WHERE tg_id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, telegram_id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Unban(ctx context.Context, telegram_id int64) error {
	query := `
		UPDATE users
		SET is_banned = FALSE
		WHERE tg_id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, telegram_id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) AddAdmin(ctx context.Context, telegram_id int64) error {
	query := `
		UPDATE users
		SET is_admin = TRUE
		WHERE tg_id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, telegram_id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) RemoveAdmin(ctx context.Context, telegram_id int64) error {
	query := `
		UPDATE users
		SET is_admin = FALSE
		WHERE tg_id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, telegram_id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT tg_id, username, first_name, last_name, is_admin, is_banned
		FROM users
		LIMIT $1 OFFSET $2
	`

	rows, err := r.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*domain.User{}

	for rows.Next() {
		var u domain.User
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.FirstName,
			&u.LastName,
			&u.IsAdmin,
			&u.IsBanned,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}
