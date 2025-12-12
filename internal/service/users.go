package service

import (
	"context"

	"github.com/Rissochek/db-cw/internal/model"
)

func (s *Service) CreateUser(ctx context.Context, user *model.User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *Service) UpdateUser(ctx context.Context, user *model.User) error {
	return s.repo.UpdateUser(ctx, user)
}

func (s *Service) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}
