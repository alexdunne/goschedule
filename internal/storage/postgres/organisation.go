package postgres

import (
	"context"
)

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
