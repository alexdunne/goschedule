package postgres

import (
	"context"
	"goschedule/internal/accounts"
	"goschedule/internal/storage"
	"time"

	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(url string) (*Storage, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	config.ConnConfig.Logger = zapadapter.NewLogger(zap.L())

	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		db: conn,
	}

	return storage, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) CreateAccount(ctx context.Context, account *accounts.NewAccount) error {
	newUser := User{
		ID:        storage.GenerateID(),
		Name:      account.Name,
		Email:     account.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newUserExternaLogin := UserExternalLogin{
		ID:        storage.GenerateID(),
		UserID:    newUser.ID,
		Source:    account.Source,
		SourceID:  account.SourceID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
	INSERT INTO users(id, name, email, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)`,
		newUser.ID,
		newUser.Name,
		newUser.Email,
		newUser.CreatedAt,
		newUser.UpdatedAt,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
	INSERT INTO user_external_logins(id, user_id, source, source_id, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`,
		newUserExternaLogin.ID,
		newUserExternaLogin.UserID,
		newUserExternaLogin.Source,
		newUserExternaLogin.SourceID,
		newUserExternaLogin.CreatedAt,
		newUserExternaLogin.UpdatedAt,
	)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	account.UserID = newUser.ID

	return nil
}
