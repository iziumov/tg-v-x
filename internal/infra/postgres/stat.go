package db

import (
	"context"
	"database/sql"
	"iziumov/tv-v-x/internal/domain"
)

type StatRepo struct {
	*DB
}

func NewStatRepo(db *DB) *StatRepo {
	return &StatRepo{db}
}

func (r *StatRepo) CreateStat(ctx context.Context, telegram_id int64) error {
	query := `
		INSERT INTO user_stats(user_id)
		VALUES ($1)
		ON CONFLICT(user_id) DO NOTHING
	`

	_, err := r.DB.ExecContext(ctx, query, telegram_id)
	if err != nil {
		return err
	}

	return nil
}

func (r *StatRepo) IncrementSuccess(ctx context.Context, telegram_id int64, bytes int64) error {
	query := `
		INSERT INTO user_stats(user_id)
		VALUES ($1)
		ON CONFLICT(user_id) DO NOTHING
	`
	_, _ = r.DB.ExecContext(ctx, query, telegram_id)

	query = `
		UPDATE user_stats
		SET total_jobs = total_jobs + 1,
		success_jobs = success_jobs + 1,
		total_bytes = total_bytes + $2,
		last_active = now()
		WHERE user_id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, telegram_id, bytes)
	if err != nil {
		return err
	}

	return nil
}

func (r *StatRepo) IncrementFailed(ctx context.Context, telegram_id int64) error {
	query := `
		UPDATE user_stats
		SET total_jobs = total_jobs + 1,
		failed_jobs = failed_jobs + 1,
		last_active = now()
		WHERE user_id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, telegram_id)
	if err != nil {
		return err
	}

	return nil
}

func (r *StatRepo) GetUserStats(ctx context.Context, telegram_id int64) (*domain.Stats, error) {
	query := `
		INSERT INTO user_stats(user_id)
		VALUES ($1)
		ON CONFLICT(user_id) DO NOTHING
	`
	_, _ = r.DB.ExecContext(ctx, query, telegram_id)

	query = `
		SELECT total_jobs, success_jobs, failed_jobs, total_bytes
		FROM user_stats
		WHERE user_id = $1
	`

	stat := &domain.Stats{}
	err := r.DB.QueryRowContext(ctx, query, telegram_id).Scan(
		&stat.TotalJobs,
		&stat.SuccessJobs,
		&stat.FailedJobs,
		&stat.TotalBytes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &domain.Stats{}, nil
		}
		return nil, err
	}

	return stat, nil
}

func (r *StatRepo) GetGlobalStats(ctx context.Context) (*domain.Stats, error) {
	query := `
		SELECT
			COALESCE(SUM(total_jobs), 0),
			COALESCE(SUM(success_jobs), 0),
			COALESCE(SUM(failed_jobs), 0),
			COALESCE(SUM(total_bytes), 0)
		FROM user_stats
	`

	stats := &domain.Stats{}
	err := r.DB.QueryRowContext(ctx, query).Scan(
		&stats.TotalJobs,
		&stats.SuccessJobs,
		&stats.FailedJobs,
		&stats.TotalBytes,
	)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
