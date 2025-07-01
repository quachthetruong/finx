package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/loanoffer/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/internal/database/mapper"
)

var _ repository.LoanPackageOfferRepository = (*LoanPackageOfferPostgresRepository)(nil)

type LoanPackageOfferPostgresRepository struct {
	getDbFunc database.GetDbFunc
}

func (r *LoanPackageOfferPostgresRepository) FindByIdWithRequest(ctx context.Context, id int64) (entity.LoanPackageOffer, error) {
	dest := LoanPackageOfferWithRequest{}
	joinTables := table.LoanPackageOffer.
		INNER_JOIN(
			table.LoanPackageRequest,
			table.LoanPackageOffer.LoanPackageRequestID.EQ(table.LoanPackageRequest.ID),
		)
	err := table.LoanPackageOffer.SELECT(
		table.LoanPackageOffer.AllColumns, table.LoanPackageRequest.AllColumns,
	).FROM(joinTables).WHERE(table.LoanPackageOffer.ID.EQ(postgres.Int64(id))).QueryContext(
		ctx, r.getDbFunc(ctx), &dest,
	)
	if err != nil {
		return entity.LoanPackageOffer{}, fmt.Errorf("LoanPackageOfferPostgresRepository FindByIdWithRequest %w", err)
	}
	return MapLoanPackageOfferWithRequestDbToEntity(dest), nil
}

func (r *LoanPackageOfferPostgresRepository) FindAllForInvestorWithRequestAndLine(ctx context.Context, filter entity.LoanPackageOfferFilter) ([]entity.LoanPackageOffer, error) {
	dest := make([]LoanPackageOfferWithRequest, 0)
	joinTables := table.LoanPackageOffer.
		INNER_JOIN(
			table.LoanPackageRequest,
			table.LoanPackageOffer.LoanPackageRequestID.EQ(table.LoanPackageRequest.ID),
		).
		LEFT_JOIN(
			table.LoanPackageOfferInterest,
			table.LoanPackageOffer.ID.EQ(table.LoanPackageOfferInterest.LoanPackageOfferID),
		).
		LEFT_JOIN(table.LoanContract, table.LoanPackageOfferInterest.ID.EQ(table.LoanContract.LoanOfferInterestID)).
		INNER_JOIN(table.Symbol, table.LoanPackageRequest.SymbolID.EQ(table.Symbol.ID))
	stm := table.LoanPackageOffer.SELECT(
		table.LoanPackageOffer.AllColumns,
		table.LoanPackageRequest.AllColumns,
		table.LoanPackageOfferInterest.AllColumns,
		table.Symbol.AllColumns,
		table.LoanContract.AllColumns,
	).FROM(joinTables).WHERE(
		ApplyFilter(filter),
	).ORDER_BY(table.LoanPackageOffer.ID.ASC(), table.LoanPackageOfferInterest.ID.DESC())
	if err := stm.QueryContext(
		ctx, r.getDbFunc(ctx), &dest,
	); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.LoanPackageOffer{}, nil
		}
		return nil, fmt.Errorf("LoanPackageOfferPostgresRepository FindAll %w", err)
	}
	return MapLoanPackageOfferWithRequestsDbToEntity(dest), nil
}

func (r *LoanPackageOfferPostgresRepository) InvestorGetById(ctx context.Context, id int64) (entity.LoanPackageOffer, error) {
	dest := LoanPackageOfferWithRequest{}
	joinTables := table.LoanPackageOffer.
		INNER_JOIN(
			table.LoanPackageRequest,
			table.LoanPackageOffer.LoanPackageRequestID.EQ(table.LoanPackageRequest.ID),
		).
		LEFT_JOIN(
			table.LoanPackageOfferInterest,
			table.LoanPackageOffer.ID.EQ(table.LoanPackageOfferInterest.LoanPackageOfferID),
		).
		LEFT_JOIN(table.LoanContract, table.LoanPackageOfferInterest.ID.EQ(table.LoanContract.LoanOfferInterestID)).
		INNER_JOIN(table.Symbol, table.LoanPackageRequest.SymbolID.EQ(table.Symbol.ID))
	err := table.LoanPackageOffer.SELECT(
		table.LoanPackageOffer.AllColumns,
		table.LoanPackageRequest.AllColumns,
		table.LoanPackageOfferInterest.AllColumns,
		table.Symbol.AllColumns,
		table.LoanContract.AllColumns,
	).FROM(joinTables).WHERE(
		table.LoanPackageOffer.ID.EQ(postgres.Int64(id)),
	).ORDER_BY(table.LoanPackageOfferInterest.ID.ASC()).
		QueryContext(ctx, r.getDbFunc(ctx), &dest)
	if err != nil {
		return entity.LoanPackageOffer{}, fmt.Errorf("LoanPackageOfferPostgresRepository InvestorGetById %w", err)
	}
	return MapLoanPackageOfferWithRequestDbToEntity(dest), nil
}

func (r *LoanPackageOfferPostgresRepository) GetExpiredOffers(ctx context.Context) ([]entity.LoanPackageOffer, error) {
	dest := make([]model.LoanPackageOffer, 0)
	err := postgres.SELECT(table.LoanPackageOffer.AllColumns).FROM(
		table.LoanPackageOffer.INNER_JOIN(
			table.LoanPackageOfferInterest,
			table.LoanPackageOfferInterest.LoanPackageOfferID.EQ(table.LoanPackageOffer.ID),
		),
	).WHERE(
		table.LoanPackageOffer.ExpiredAt.LT(postgres.TimestampT(time.Now())).
			AND(table.LoanPackageOffer.FlowType.EQ(postgres.String(entity.FlowTypeDnseOnline.String()))).
			AND(table.LoanPackageOfferInterest.Status.EQ(postgres.String(entity.LoanPackageOfferInterestStatusPending.String()))),
	).QueryContext(
		ctx, r.getDbFunc(ctx), &dest,
	)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.LoanPackageOffer{}, nil
		}
		return nil, fmt.Errorf("LoanPackageOfferPostgresRepository GetExpired %w", err)
	}
	return mapper.MapLoanPackageOffersDbToEntity(dest), nil
}

func (r *LoanPackageOfferPostgresRepository) Create(ctx context.Context, loanPackageOffer entity.LoanPackageOffer) (entity.LoanPackageOffer, error) {
	createModel := MapLoanPackageOfferEntityToDb(loanPackageOffer)
	created := model.LoanPackageOffer{}
	if err := table.LoanPackageOffer.INSERT(table.LoanPackageOffer.MutableColumns).MODEL(createModel).RETURNING(table.LoanPackageOffer.AllColumns).QueryContext(
		ctx, r.getDbFunc(ctx), &created,
	); err != nil {
		return entity.LoanPackageOffer{}, fmt.Errorf("LoanPackageOfferPostgresRepository Create %w", err)
	}
	return mapper.MapLoanPackageOfferDbToEntity(created), nil
}

func (r *LoanPackageOfferPostgresRepository) BulkCreate(ctx context.Context, loanPackageOffers []entity.LoanPackageOffer) ([]entity.LoanPackageOffer, error) {
	insertModels := MapLoanPackageOffersEntityToDb(loanPackageOffers)
	res := make([]model.LoanPackageOffer, 0, len(loanPackageOffers))
	err := table.LoanPackageOffer.INSERT(table.LoanPackageOffer.MutableColumns).
		MODELS(insertModels).
		RETURNING(table.LoanPackageOffer.AllColumns).
		QueryContext(ctx, r.getDbFunc(ctx), &res)
	if err != nil {
		return []entity.LoanPackageOffer{}, fmt.Errorf("LoanPackageOfferPostgresRepository BulkCreate %w", err)
	}
	return MapLoanPackageOffersDbToEntity(res), nil
}

func NewLoanPackageOfferPostgresRepository(getDbFunc database.GetDbFunc) *LoanPackageOfferPostgresRepository {
	return &LoanPackageOfferPostgresRepository{getDbFunc: getDbFunc}
}
