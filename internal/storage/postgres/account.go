package postgres

import (
	"context"
	"goschedule/internal/accounts"
	"goschedule/internal/storage"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

func (s *Storage) CreateAccount(ctx context.Context, account *accounts.Account) error {
	tx, err := s.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// First check if we already have an existing auth for this source and source id combo
	if existing, err := findUserExternalLoginBySourceID(
		ctx,
		tx,
		account.Source,
		account.SourceID,
	); err == nil {
		zap.S().Debugf("UserExternalLogin aleady exists for %+v", account)
		user, err := findUserByID(ctx, tx, existing.UserID)

		if err != nil {
			return err
		}

		account.UserID = user.ID
		return tx.Commit(ctx)
	} else if err != pgx.ErrNoRows {
		return err
	}

	zap.S().Debugf("UserExternalLogin does not exist for %+v. Creating a new user and external login now", account)

	user := User{
		ID:    storage.GenerateID(),
		Name:  account.Name,
		Email: account.Email,
	}

	userExternaLogin := UserExternalLogin{
		ID:       storage.GenerateID(),
		UserID:   user.ID,
		Source:   account.Source,
		SourceID: account.SourceID,
	}

	org := &Organisation{
		ID:      storage.GenerateID(),
		Name:    "My Org",
		OwnerID: user.ID,
	}

	schedule := &Schedule{
		ID:             storage.GenerateID(),
		Name:           "My Schedule",
		OwnerID:        user.ID,
		OrganisationID: org.ID,
	}

	if err := createUser(ctx, tx, &user); err != nil {
		return err
	}

	if err := createUserExternalLogin(ctx, tx, &userExternaLogin); err != nil {
		return err
	}

	if err := createOrganisation(ctx, tx, org); err != nil {
		return err
	}

	if err := createSchedule(ctx, tx, schedule); err != nil {
		return err
	}

	account.UserID = user.ID

	return tx.Commit(ctx)
}
