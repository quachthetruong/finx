package dbtest

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"

	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
)

var Container *gnomock.Container

const (
	defaultPassword = "Encap@1234"
	defaultUser     = "encapital"
	dbPrefix        = "finoffer"
)

func NewDb(t *testing.T) (newDb *sql.DB, dropDbFunc func(), truncateFunc func()) {
	if Container == nil {
		t.Fatal("Container is not yet started")
	}
	db, err := connectBaseDb()
	if err != nil {
		t.Fatal("err connect to base db", slog.String("err", err.Error()))
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal("err close connection to db", slog.String("err", err.Error()))
		}
	}()
	dbName := fmt.Sprintf("%s_%d", dbPrefix, rand.Int63())
	if err := createDatabase(db, dbName); err != nil {
		slog.Error("err creating postgres database", slog.String("err", err.Error()))
		os.Exit(1)
	}
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable", "encapital", "Encap@1234", Container.Host,
		Container.DefaultPort(), dbName,
	)
	sub, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("err connect postgres database", slog.String("err", err.Error()))
		os.Exit(1)
	}
	// apply migrations
	if err := database.MigrationUp(connStr); err != nil {
		slog.Error("err connect postgres database", slog.String("err", err.Error()))
		os.Exit(1)
	}
	return sub, func() {
			_ = deleteDatabase(dbName)
		}, func() {
			truncateAll(sub)
		}
}

func InitPostgres() (cleanUp func()) {
	slog.Info("spinning up postgres container for testing..")
	p := postgres.Preset(
		postgres.WithUser(defaultUser, defaultPassword),
		postgres.WithDatabase(dbPrefix),
	)
	container, err := gnomock.Start(p, gnomock.WithTimeout(60*time.Second), gnomock.WithDebugMode())
	if err != nil {
		slog.Error("err creating postgres container", slog.String("err", err.Error()))
		if err := gnomock.Stop(container); err != nil {
			slog.Error("err stopping postgres container", slog.String("err", err.Error()))
		}
		os.Exit(1)
	}
	Container = container
	return func() {
		if err := gnomock.Stop(container); err != nil {
			slog.Error("err stopping postgres container", slog.String("err", err.Error()))
		}
		slog.Info("Clean database test container success")
	}
}

func connectBaseDb() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Container.Host, Container.DefaultPort(),
		defaultUser, defaultPassword, dbPrefix,
	)
	return sql.Open("postgres", connStr)
}

func createDatabase(db *sql.DB, name string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", name))

	return err
}

func deleteDatabase(name string) error {
	db, err := connectBaseDb()
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s;", name))
	return err
}

func truncateAll(db database.DB) {
	tables := []string{
		table.Symbol.TableName(),
		table.SymbolScore.TableName(),
		table.StockExchange.TableName(),
		table.LoanPackageRequest.TableName(),
		table.LoanPackageOffer.TableName(),
		table.LoanPackageOfferInterest.TableName(),
		table.LoanContract.TableName(),
		table.ScoreGroup.TableName(),
		table.ScoreGroupInterest.TableName(),
		table.BlacklistSymbol.TableName(),
		table.Investor.TableName(),
		table.LoanPolicyTemplate.TableName(),
		table.InvestorAccount.TableName(),
		table.SubmissionSheetMetadata.TableName(),
		table.SubmissionSheetDetail.TableName(),
		table.SuggestedOfferConfig.TableName(),
		table.FinancialConfiguration.TableName(),
		table.SuggestedOffer.TableName(),
		table.PromotionCampaign.TableName(),
	}
	for _, t := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE %s RESTART IDENTITY CASCADE;", t))
		if err != nil {
			slog.Error("err truncate table", slog.String("err", err.Error()))
		}
	}
}
