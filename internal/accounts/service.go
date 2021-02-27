package accounts

import (
	"context"
)

type Repository interface {
	CreateAccount(ctx context.Context, account *Account) error
	CreateOrganisation(ctx context.Context, organisation *Organisation) error
	CreateSchedule(ctx context.Context, schedule *Schedule) error
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

	if err := s.repo.CreateOrganisation(ctx, &Organisation{Name: "My Organisation", OwnerID: account.UserID}); err != nil {
		return err
	}

	if err := s.repo.CreateSchedule(ctx, &Schedule{Name: "Default Schedule", OwnerID: account.UserID}); err != nil {
		return err
	}

	return nil
}
