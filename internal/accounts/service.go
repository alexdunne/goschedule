package accounts

import (
	"context"
)

type Repository interface {
	CreateAccount(ctx context.Context, account *Account) error
}

type Service interface {
	CreateAccount(ctx context.Context, account *Account) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	s := &service{
		repo: repo,
	}

	return s
}

func (s *service) CreateAccount(ctx context.Context, account *Account) error {
	if err := account.Validate(); err != nil {
		return err
	}

	if err := s.repo.CreateAccount(ctx, account); err != nil {
		return err
	}

	return nil
}
