package main

import (
	"context"
	"errors"

	"github.com/lozhkindm/otus-banners-rotation/internal/app"
	sqlstorage "github.com/lozhkindm/otus-banners-rotation/internal/storage/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const storageTypeSQL = "sql"

var undefinedStorageType = errors.New("undefined storage type")

func NewStorage(ctx context.Context, config Config) (app.Storage, func(context.Context) error, error) {
	var (
		storage   app.Storage
		closeFunc func(ctx2 context.Context) error
	)

	switch config.App.StorageType {
	case storageTypeSQL:
		db, err := sqlx.Open("pgx", config.PostgreSQL.BuildDSN())
		if err != nil {
			return nil, nil, err
		}
		ss := sqlstorage.New(db)
		if err := ss.Connect(ctx); err != nil {
			return nil, nil, err
		}
		storage = ss
		closeFunc = ss.Close
	default:
		return nil, nil, undefinedStorageType
	}

	return storage, closeFunc, nil
}
