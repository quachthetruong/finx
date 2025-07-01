package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/loanofferinterest/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/pkg/querymod"
)

var _ repository.LoanPackageOfferInterestRepository = &LoanPackageOfferInterestPostgresRepository{}

type LoanPackageOfferInterestPostgresRepository struct {
	getDbFunc database.GetDbFunc
}

func (l *LoanPackageOfferInterestPostgresRepository) GetByOfferIdWithLock(ctx context.Context, offerId int64) ([]entity.LoanPackageOfferInterest, error) {
	loanPackageOfferInterests := make([]model.LoanPackageOfferInterest, 0)
	stm := table.LoanPackageOfferInterest.
		SELECT(table.LoanPackageOfferInterest.AllColumns).
		WHERE(table.LoanPackageOfferInterest.LoanPackageOfferID.EQ(postgres.Int64(offerId))).
		FOR(postgres.UPDATE().NOWAIT())
	if err := stm.QueryContext(ctx, l.getDbFunc(ctx), &loanPackageOfferInterests); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.LoanPackageOfferInterest{}, nil
		}
		return nil, fmt.Errorf(
			"LoanPackageOfferInterestPostgresRepository GetByOfferIdWithLock %w", err,
		)
	}
	return MapLoanPackageOfferInterestsDbToEntity(loanPackageOfferInterests), nil
}

func (l *LoanPackageOfferInterestPostgresRepository) Create(ctx context.Context, loanPackageOfferInterest entity.LoanPackageOfferInterest) (entity.LoanPackageOfferInterest, error) {
	created := model.LoanPackageOfferInterest{}
	err := table.LoanPackageOfferInterest.
		INSERT(table.LoanPackageOfferInterest.MutableColumns).
		MODEL(MapLoanPackageOfferInterestEntityToDb(loanPackageOfferInterest)).
		RETURNING(table.LoanPackageOfferInterest.AllColumns).
		QueryContext(
			ctx, l.getDbFunc(ctx), &created,
		)
	if err != nil {
		return entity.LoanPackageOfferInterest{}, fmt.Errorf(
			"LoanPackageOfferInterestPostgresRepository Create %w", err,
		)
	}
	return MapLoanPackageOfferInterestDbToEntity(created), nil
}

func (l *LoanPackageOfferInterestPostgresRepository) CountWithFilter(ctx context.Context, filter entity.OfferInterestFilter) (int64, error) {
	dest := struct {
		Count int64
	}{}
	stm := table.LoanPackageOfferInterest.SELECT(postgres.COUNT(table.LoanPackageOfferInterest.ID)).WHERE(ApplyFilter(filter))
	if err := stm.QueryContext(ctx, l.getDbFunc(ctx), &dest); err != nil {
		return 0, fmt.Errorf("LoanPackageOfferInterestPostgresRepository CountWithFilter %w", err)
	}
	return dest.Count, nil
}

func (l *LoanPackageOfferInterestPostgresRepository) Update(ctx context.Context, offerInterest entity.LoanPackageOfferInterest) (entity.LoanPackageOfferInterest, error) {
	updated := model.LoanPackageOfferInterest{}
	err := table.LoanPackageOfferInterest.
		UPDATE(table.LoanPackageOfferInterest.MutableColumns).
		MODEL(MapLoanPackageOfferInterestEntityToDb(offerInterest)).
		WHERE(table.LoanPackageOfferInterest.ID.EQ(postgres.Int64(offerInterest.Id))).
		RETURNING(table.LoanPackageOfferInterest.AllColumns).
		QueryContext(
			ctx, l.getDbFunc(ctx), &updated,
		)
	if err != nil {
		return entity.LoanPackageOfferInterest{}, fmt.Errorf(
			"LoanPackageOfferInterestPostgresRepository Update %w", err,
		)
	}
	return MapLoanPackageOfferInterestDbToEntity(updated), nil
}

func (l *LoanPackageOfferInterestPostgresRepository) GetWithFilter(ctx context.Context, filter entity.OfferInterestFilter) ([]entity.LoanPackageOfferInterest, error) {
	stm := table.LoanPackageOfferInterest.
		SELECT(
			table.LoanPackageOfferInterest.AllColumns, table.LoanContract.AllColumns, table.LoanPackageOffer.AllColumns,
		).FROM(
		table.LoanPackageOfferInterest.LEFT_JOIN(
			table.LoanContract, table.LoanPackageOfferInterest.ID.EQ(table.LoanContract.LoanOfferInterestID),
		).INNER_JOIN(
			table.LoanPackageOffer, table.LoanPackageOfferInterest.LoanPackageOfferID.EQ(table.LoanPackageOffer.ID),
		),
	).WHERE(ApplyFilter(filter))
	if orderClause := ApplySort(filter); len(orderClause) > 0 {
		stm = stm.ORDER_BY(orderClause...)
	}
	if limit := filter.Limit(); limit > 0 {
		stm = stm.LIMIT(limit).OFFSET(filter.Offset())
	}
	loanPackageOfferInterests := make([]OfferInterestWithContract, 0)
	if err := stm.QueryContext(
		ctx, l.getDbFunc(ctx), &loanPackageOfferInterests,
	); err != nil {
		return nil, fmt.Errorf("LoanPackageOfferInterestPostgresRepository GetWithFilter %w", err)
	}
	return MapOfferInterestWithContractsDbToEntity(loanPackageOfferInterests), nil
}

func (l *LoanPackageOfferInterestPostgresRepository) UpdateStatus(ctx context.Context, ids []int64, status entity.LoanPackageOfferInterestStatus) error {
	if _, err := table.LoanPackageOfferInterest.
		UPDATE(table.LoanPackageOfferInterest.Status).
		SET(postgres.String(status.String())).
		WHERE(table.LoanPackageOfferInterest.ID.IN(querymod.In(ids)...)).
		ExecContext(ctx, l.getDbFunc(ctx)); err != nil {
		return fmt.Errorf("LoanPackageOfferInterestPostgresRepository UpdateStatus %w", err)
	}
	return nil
}

func (l *LoanPackageOfferInterestPostgresRepository) GetById(ctx context.Context, id int64, opts ...querymod.GetOption) (entity.LoanPackageOfferInterest, error) {
	loanPackageOfferInterest := model.LoanPackageOfferInterest{}
	getQm := querymod.GetQm{}
	for _, opt := range opts {
		opt(&getQm)
	}
	stm := table.LoanPackageOfferInterest.
		SELECT(table.LoanPackageOfferInterest.AllColumns).
		WHERE(table.LoanPackageOfferInterest.ID.EQ(postgres.Int64(id)))
	if getQm.ForUpdate {
		stm = stm.FOR(postgres.UPDATE().NOWAIT())
	}
	if err := stm.QueryContext(ctx, l.getDbFunc(ctx), &loanPackageOfferInterest); err != nil {
		return entity.LoanPackageOfferInterest{}, fmt.Errorf(
			"LoanPackageOfferInterestPostgresRepository GetById %w", err,
		)
	}
	return MapLoanPackageOfferInterestDbToEntity(loanPackageOfferInterest), nil
}

func (l *LoanPackageOfferInterestPostgresRepository) GetByIds(ctx context.Context, ids []int64, opts ...querymod.GetOption) ([]entity.LoanPackageOfferInterest, error) {
	loanPackageOfferInterest := make([]model.LoanPackageOfferInterest, 0, len(ids))
	getQm := querymod.GetQm{}
	for _, opt := range opts {
		opt(&getQm)
	}
	stm := table.LoanPackageOfferInterest.
		SELECT(table.LoanPackageOfferInterest.AllColumns).
		WHERE(table.LoanPackageOfferInterest.ID.IN(querymod.In(ids)...))
	if getQm.ForUpdate {
		stm = stm.FOR(postgres.UPDATE().NOWAIT())
	}
	if err := stm.QueryContext(ctx, l.getDbFunc(ctx), &loanPackageOfferInterest); err != nil {
		return nil, fmt.Errorf(
			"LoanPackageOfferInterestPostgresRepository GetById %w", err,
		)
	}
	if len(loanPackageOfferInterest) != len(ids) {
		return nil, apperrors.ErrorLoanPackageOfferInterestNotFound
	}
	return MapLoanPackageOfferInterestsDbToEntity(loanPackageOfferInterest), nil
}

func (l *LoanPackageOfferInterestPostgresRepository) BulkCreate(ctx context.Context, loanPackageOfferInterests []entity.LoanPackageOfferInterest) ([]entity.LoanPackageOfferInterest, error) {
	createModels := make([]model.LoanPackageOfferInterest, 0, len(loanPackageOfferInterests))
	for _, loanPackageOfferInterest := range loanPackageOfferInterests {
		createModels = append(createModels, MapLoanPackageOfferInterestEntityToDb(loanPackageOfferInterest))
	}
	res := make([]model.LoanPackageOfferInterest, 0, len(loanPackageOfferInterests))
	if err := table.LoanPackageOfferInterest.
		INSERT(table.LoanPackageOfferInterest.MutableColumns).
		MODELS(createModels).
		RETURNING(table.LoanPackageOfferInterest.AllColumns).
		QueryContext(ctx, l.getDbFunc(ctx), &res); err != nil {
		return nil, fmt.Errorf("LoanPackageOfferInterestPostgresRepository BulkCreate %w", err)
	}
	return MapLoanPackageOfferInterestsDbToEntity(res), nil
}

func (l *LoanPackageOfferInterestPostgresRepository) CancelByOfferId(ctx context.Context, offerId int64, cancelledBy string, cancelledReason entity.CancelledReason) error {
	if _, err := table.LoanPackageOfferInterest.UPDATE().
		SET(
			table.LoanPackageOfferInterest.Status.SET(postgres.String(entity.LoanPackageOfferInterestStatusCancelled.String())),
			table.LoanPackageOfferInterest.CancelledBy.SET(postgres.String(cancelledBy)),
			table.LoanPackageOfferInterest.CancelledAt.SET(postgres.TimestampT(time.Now())),
			table.LoanPackageOfferInterest.CancelledReason.SET(postgres.String(cancelledReason.String())),
		).
		WHERE(
			table.LoanPackageOfferInterest.LoanPackageOfferID.EQ(postgres.Int64(offerId)).
				AND(table.LoanPackageOfferInterest.Status.EQ(postgres.String(entity.LoanPackageOfferInterestStatusPending.String()))),
		).
		ExecContext(ctx, l.getDbFunc(ctx)); err != nil {
		return fmt.Errorf("LoanPackageOfferInterestPostgresRepository CancelByOfferId %w", err)
	}
	return nil
}

func (l *LoanPackageOfferInterestPostgresRepository) CancelExpiredOfferInterests(ctx context.Context, offerIds []int64) error {
	if len(offerIds) == 0 {
		return nil
	}
	if _, err := table.LoanPackageOfferInterest.UPDATE().
		SET(
			table.LoanPackageOfferInterest.Status.SET(postgres.String(entity.LoanPackageOfferInterestStatusCancelled.String())),
			table.LoanPackageOfferInterest.CancelledBy.SET(postgres.String("system")),
			table.LoanPackageOfferInterest.CancelledAt.SET(postgres.TimestampT(time.Now())),
			table.LoanPackageOfferInterest.CancelledReason.SET(postgres.String(entity.LoanPackageOfferCancelledReasonExpired.String())),
		).WHERE(
		table.LoanPackageOfferInterest.LoanPackageOfferID.IN(querymod.In(offerIds)...).
			AND(table.LoanPackageOfferInterest.Status.EQ(postgres.String(entity.LoanPackageOfferInterestStatusPending.String()))),
	).ExecContext(ctx, l.getDbFunc(ctx)); err != nil {
		return fmt.Errorf("LoanPackageOfferInterestPostgresRepository CancelExpiredOfferInterests %w", err)
	}
	return nil
}

func (l *LoanPackageOfferInterestPostgresRepository) GetRequestBasedLoanOfferInterests(ctx context.Context) ([]entity.LoanPackageOfferInterest, error) {
	dest := make([]model.LoanPackageOfferInterest, 0)
	err := table.LoanPackageOfferInterest.
		SELECT(table.LoanPackageOfferInterest.AllColumns).
		WHERE(
			table.LoanPackageOfferInterest.LoanID.NOT_EQ(postgres.Int64(0)).
				AND(table.LoanPackageOfferInterest.AssetType.EQ(postgres.NewEnumValue(entity.AssetTypeUnderlying.String()))),
		).QueryContext(
		ctx, l.getDbFunc(ctx), &dest,
	)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.LoanPackageOfferInterest{}, nil
		}
		return nil, fmt.Errorf("LoanPackageRequestPostgresRepository GetRequestBasedLoanOfferInterests: %w", err)
	}
	return MapLoanPackageOfferInterestsDbToEntity(dest), nil
}

func NewLoanPackageOfferInterestPostgresRepository(getDbFunc database.GetDbFunc) *LoanPackageOfferInterestPostgresRepository {
	return &LoanPackageOfferInterestPostgresRepository{
		getDbFunc: getDbFunc,
	}
}
