package db

import (
	"context"
	"iziumov/tv-v-x/internal/domain"
	"iziumov/tv-v-x/internal/dto"
)

type JobRepo struct {
	*DB
}

func NewJobRepo(db *DB) *JobRepo {
	return &JobRepo{db}
}

func (r *JobRepo) CreateJob(ctx context.Context, job dto.CreateJob) (int, error) {
	query := `
		INSERT INTO jobs(user_id, url, status, platform)
		VALUES ($1, $2, 'pending', $3)
		RETURNING id
	`

	var id int
	err := r.QueryRowContext(ctx, query, job.TgId, job.Url, job.Platform).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *JobRepo) UpdateFinalStatus(ctx context.Context, job dto.UpdateJob) error {
	query := `
		UPDATE jobs
		SET 
			status = $2, 
			error_msg = COALESCE($3, error_msg), 
			finished_at = now()
		WHERE id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, job.ID, job.Status, job.ErrorMsg)
	if err != nil {
		return err
	}

	return nil
}

func (r *JobRepo) GetFailedJobs(ctx context.Context) ([]*domain.Job, error) {
	query := `
		SELECT id, user_id, url, status, platform, 
			COALESCE(file_size_bytes, 0), 
			COALESCE(error_msg, ""), 
			created_at, 
			COALESCE(finished_at, created_at)
		FROM jobs
		WHERE status = 'failed'
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*domain.Job
	for rows.Next() {
		j := &domain.Job{}
		if err := rows.Scan(
			&j.Id, &j.UserId, &j.Url, &j.Status, &j.Platform,
			&j.FileSize, &j.ErrorMsg, &j.CreatedAt, &j.FinishedAt,
		); err != nil {
			return nil, err
		}

		jobs = append(jobs, j)
	}

	return jobs, nil
}

func (r *JobRepo) GetHistoryByTgId(ctx context.Context, telegram_id int64) ([]*domain.JobRecord, error) {
	query := `
		SELECT url
		FROM jobs
		WHERE user_id = $1
	`

	rows, err := r.DB.QueryContext(ctx, query, telegram_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*domain.JobRecord
	for rows.Next() {
		j := &domain.JobRecord{}
		if err := rows.Scan(&j.Url); err != nil {
			return nil, err
		}

		jobs = append(jobs, j)
	}

	return jobs, nil
}
