package atomicity

import (
	"context"
	"database/sql"

	"financing-offer/internal/apperrors"
)

const TxKey = "transactionInstance"

type DbAtomicExecutor struct {
	DB *sql.DB
}

func (e *DbAtomicExecutor) Execute(parentCtx context.Context, executeFunc func(ctx context.Context) error) (err error) {
	tx, err := e.DB.Begin()
	if err != nil {
		return apperrors.New(err, apperrors.WithCode(500), apperrors.WithMessage("begin transaction"))
	}
	transactionalCtx := ContextSetTx(parentCtx, tx)
	defer func() {
		if r := recover(); r != nil {
			if scopedErr := tx.Rollback(); scopedErr != nil {
				err = scopedErr
			}
			panic(r)
		}
		if err != nil {
			if scopedErr := tx.Rollback(); scopedErr != nil {
				err = scopedErr
			}
		} else {
			if scopedErr := tx.Commit(); scopedErr != nil {
				err = scopedErr
			}
		}
	}()
	err = executeFunc(transactionalCtx)
	return err
}

func ContextSetTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func WithIgnoreTx(cte context.Context) context.Context {
	return context.WithValue(cte, TxKey, nil)
}

func ContextGetTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}
