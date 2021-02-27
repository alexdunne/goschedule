package accounts

import "context"

type Repository interface {
	CreateAccount(ctx context.Context, account *NewAccount) error
}

type Service interface {
	CreateAccount(ctx context.Context, account *NewAccount) error
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

func (s *service) CreateAccount(ctx context.Context, account *NewAccount) error {
	if err := account.Validate(); err != nil {
		return err
	}

	s.repo.CreateAccount(ctx, account)
	return nil
}
