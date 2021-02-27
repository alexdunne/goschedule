package postgres

import (
	"context"
	"goschedule/internal/accounts"
	"goschedule/internal/storage"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	db *pgxpool.Pool
}

type Tx struct {
	pgx.Tx
	db  *pgxpool.Pool
	now time.Time
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

func (s *Storage) BeginTx(ctx context.Context) (*Tx, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &Tx{
		Tx:  tx,
		db:  s.db,
		now: time.Now(),
	}, nil
}

func (s *Storage) CreateAccount(ctx context.Context, account *accounts.NewAccount) error {
	newUser := User{
		ID:    storage.GenerateID(),
		Name:  account.Name,
		Email: account.Email,
	}

	newUserExternaLogin := UserExternalLogin{
		ID:       storage.GenerateID(),
		UserID:   newUser.ID,
		Source:   account.Source,
		SourceID: account.SourceID,
	}

	tx, err := s.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := createUser(ctx, tx, &newUser); err != nil {
		return err
	}

	if err := createUserExternalLogin(ctx, tx, &newUserExternaLogin); err != nil {
		return err
	}

	account.UserID = newUser.ID

	return tx.Commit(ctx)
}

func createUser(ctx context.Context, tx *Tx, user *User) error {
	user.CreatedAt = tx.now
	user.UpdatedAt = tx.now

	_, err := tx.Exec(ctx, `
	INSERT INTO users(id, name, email, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)`,
		user.ID,
		user.Name,
		user.Email,
	)
	if err != nil {
		return err
	}

	return nil
}

func createUserExternalLogin(ctx context.Context, tx pgx.Tx, userExternalLogin *UserExternalLogin) error {
	_, err := tx.Exec(ctx, `
	INSERT INTO user_external_logins(id, user_id, source, source_id, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`,
		userExternalLogin.ID,
		userExternalLogin.UserID,
		userExternalLogin.Source,
		userExternalLogin.SourceID,
		userExternalLogin.CreatedAt,
		userExternalLogin.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
