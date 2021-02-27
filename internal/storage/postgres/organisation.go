package postgres

import (
	"context"
	"goschedule/internal/accounts"
	"goschedule/internal/storage"
)

func (s *Storage) CreateOrganisation(ctx context.Context, organisation *accounts.Organisation) error {
	tx, err := s.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	org := &Organisation{
		ID:      storage.GenerateID(),
		Name:    organisation.Name,
		OwnerID: organisation.OwnerID,
	}

	if err := createOrganisation(ctx, tx, org); err != nil {
		return nil
	}

	organisation.ID = org.ID

	return tx.Commit(ctx)
}

func createOrganisation(ctx context.Context, tx *Tx, org *Organisation) error {
	org.CreatedAt = tx.now
	org.UpdatedAt = tx.now

	_, err := tx.Exec(ctx, `
	INSERT INTO organisations(id, name, owner_id, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)`,
		org.ID,
		org.Name,
		org.OwnerID,
		org.CreatedAt,
		org.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
