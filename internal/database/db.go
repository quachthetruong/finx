package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"

	"financing-offer/assets"
	"financing-offer/internal/atomicity"
	"financing-offer/internal/config"
	"financing-offer/pkg/shutdown"
)

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type (
	GetDbFunc func(ctx context.Context) DB
)

func New(cfg config.DbConfig, tasks *shutdown.Tasks) (GetDbFunc, *atomicity.DbAtomicExecutor, error) {
	emptyAtomicExecutor := &atomicity.DbAtomicExecutor{}
	completeDsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?binary_parameters=yes", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName,
	)
	if !cfg.EnableSsl {
		completeDsn += "&sslmode=disable"
	}
	db, err := sql.Open("postgres", completeDsn)
	if err != nil {
		return nil, emptyAtomicExecutor, err
	}
	err = db.Ping()
	if err != nil {
		return nil, emptyAtomicExecutor, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	if cfg.AutoMigrate {
		err = MigrationUp(completeDsn)
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		case err != nil:
			return nil, emptyAtomicExecutor, err
		}
	}
	getDbFunc := func(ctx context.Context) DB {
		if tx := atomicity.ContextGetTx(ctx); tx != nil {
			return tx
		}
		return db
	}

	tasks.AddShutdownTask(
		func(_ context.Context) error {
			return db.Close()
		},
	)

	return getDbFunc, &atomicity.DbAtomicExecutor{DB: db}, nil
}

func MigrationUp(completeDsn string) error {
	iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, completeDsn)
	if err != nil {
		return err
	}

	return migrator.Up()
}
