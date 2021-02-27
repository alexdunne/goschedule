package postgres

import (
	"context"
)

func createSchedule(ctx context.Context, tx *Tx, schedule *Schedule) error {
	schedule.CreatedAt = tx.now
	schedule.UpdatedAt = tx.now

	_, err := tx.Exec(ctx, `
	INSERT INTO schedules(id, name, owner_id, organisation_id, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`,
		schedule.ID,
		schedule.Name,
		schedule.OwnerID,
		schedule.OrganisationID,
		schedule.CreatedAt,
		schedule.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
