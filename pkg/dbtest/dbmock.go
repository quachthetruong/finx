package dbtest

import (
	"github.com/DATA-DOG/go-sqlmock"

	"financing-offer/internal/database"
)

func New() (database.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		return nil, nil, err
	}
	return db, mock, err
}
