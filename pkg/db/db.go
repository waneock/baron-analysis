package db

import (
	"context"
	"database/sql"
	"time"
)

type DBConfigData struct {
	ConnUrl         string
	ConnTimeout     time.Duration
	ConnMaxIdleTime time.Duration
	ConnMaxLifeTime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

func New(cfg DBConfigData) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.ConnUrl)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetConnMaxLifetime(cfg.ConnMaxLifeTime)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnTimeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
