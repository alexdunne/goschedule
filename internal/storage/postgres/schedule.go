package postgres

import (
	"context"
	"goschedule/internal/accounts"
	"goschedule/internal/storage"
)

func (s *Storage) CreateSchedule(ctx context.Context, schedule *accounts.Schedule) error {
	tx, err := s.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sch := &Schedule{
		ID:      storage.GenerateID(),
		Name:    schedule.Name,
		OwnerID: schedule.OwnerID,
	}

	if err := createSchedule(ctx, tx, sch); err != nil {
		return nil
	}

	schedule.ID = sch.ID

	return tx.Commit(ctx)
}

func createSchedule(ctx context.Context, tx *Tx, schedule *Schedule) error {
	schedule.CreatedAt = tx.now
	schedule.UpdatedAt = tx.now

	_, err := tx.Exec(ctx, `
	INSERT INTO schedules(id, name, owner_id, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)`,
		schedule.ID,
		schedule.Name,
		schedule.OwnerID,
		schedule.CreatedAt,
		schedule.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
