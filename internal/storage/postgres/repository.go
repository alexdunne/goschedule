package postgres

import (
	"context"
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
