package atomicity

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAtomicExecutor(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Run(
		"tx commit", func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectCommit()
			executor := DbAtomicExecutor{DB: db}
			err := executor.Execute(
				context.Background(), func(tc context.Context) error {
					if tx := ContextGetTx(tc); tx == nil {
						t.Error("tx not exist")
					}
					return nil
				},
			)
			if err != nil {
				t.Errorf("%v", err)
			}
		},
	)
	t.Run(
		"cannot begin tx", func(t *testing.T) {
			mock.ExpectBegin().WillReturnError(errors.New("begin error"))
			executor := DbAtomicExecutor{DB: db}
			err := executor.Execute(
				context.Background(), func(tc context.Context) error {
					return nil
				},
			)
			assert.NotNil(t, err)
		},
	)
	t.Run(
		"tx rollback", func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectRollback()
			executor := DbAtomicExecutor{DB: db}
			err := executor.Execute(
				context.Background(), func(tc context.Context) error {
					return errors.New("expected")
				},
			)
			assert.NotNil(t, err)
		},
	)

	t.Run(
		"tx rollback", func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectRollback().WillReturnError(errors.New("error rollback"))
			executor := DbAtomicExecutor{DB: db}
			err := executor.Execute(
				context.Background(), func(tc context.Context) error {
					return errors.New("expected")
				},
			)
			assert.NotNil(t, err)
		},
	)

	t.Run(
		"panic during tx", func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectRollback().WillReturnError(errors.New("error rollback"))
			executor := DbAtomicExecutor{DB: db}

			assert.Panicsf(
				t, func() {
					_ = executor.Execute(
						context.Background(), func(tc context.Context) error {
							panic("expected")
						},
					)
				}, "",
			)
		},
	)

	t.Run(
		"tx commit error", func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			executor := DbAtomicExecutor{DB: db}
			err := executor.Execute(
				context.Background(), func(tc context.Context) error {
					return nil
				},
			)
			assert.NotNil(t, err)
		},
	)
}
