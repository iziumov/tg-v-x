package service

import (
	"context"
	"iziumov/tv-v-x/internal/domain"
	"iziumov/tv-v-x/internal/dto"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetByTelegramID(ctx context.Context, telegram_id int64) (*domain.User, error) {
	return s.repo.GetByTgId(ctx, telegram_id)
}

func (s *UserService) Register(ctx context.Context, user dto.CreateUser) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) Ban(ctx context.Context, telegram_id int64) error {
	return s.repo.Ban(ctx, telegram_id)
}

func (s *UserService) Unban(ctx context.Context, telegram_id int64) error {
	return s.repo.Unban(ctx, telegram_id)
}

func (s *UserService) AddAdmin(ctx context.Context, telegram_id int64) error {
	return s.repo.AddAdmin(ctx, telegram_id)
}

func (s *UserService) RemoveAdmin(ctx context.Context, telegram_id int64) error {
	return s.repo.RemoveAdmin(ctx, telegram_id)
}

func (s *UserService) GetAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return s.repo.GetAll(ctx, limit, offset)
}
