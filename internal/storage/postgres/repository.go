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

	if err := createUser(ctx, tx, &user); err != nil {
		return err
	}

	if err := createUserExternalLogin(ctx, tx, &userExternaLogin); err != nil {
		return err
	}

	account.UserID = user.ID

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
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func createUserExternalLogin(ctx context.Context, tx *Tx, userExternalLogin *UserExternalLogin) error {
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

func findUserByID(ctx context.Context, tx *Tx, id string) (*User, error) {
	sqlStmt := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := tx.QueryRow(ctx, sqlStmt, id)

	var user User
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func findUserExternalLoginBySourceID(ctx context.Context, tx *Tx, source, sourceID string) (*UserExternalLogin, error) {
	sqlStmt := `
		SELECT id, user_id, source, source_id, created_at, updated_at
		FROM user_external_logins
		WHERE source = $1 AND source_id = $2
	`

	row := tx.QueryRow(ctx, sqlStmt, source, sourceID)

	var externalLogin UserExternalLogin
	if err := row.Scan(
		&externalLogin.ID,
		&externalLogin.UserID,
		&externalLogin.Source,
		&externalLogin.SourceID,
		&externalLogin.CreatedAt,
		&externalLogin.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &externalLogin, nil
}
