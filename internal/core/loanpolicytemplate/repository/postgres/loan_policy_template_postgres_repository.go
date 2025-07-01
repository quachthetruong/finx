package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/loanpolicytemplate/repository"
	"financing-offer/internal/database"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
	"financing-offer/internal/database/dbmodels/finoffer/public/table"
	"financing-offer/pkg/querymod"
)

var _ repository.LoanPolicyTemplateRepository = (*LoanPolicyTemplateRepository)(nil)

type LoanPolicyTemplateRepository struct {
	getDbFunc database.GetDbFunc
}

func NewLoanPolicyTemplateRepository(getDbFunc database.GetDbFunc) *LoanPolicyTemplateRepository {
	return &LoanPolicyTemplateRepository{
		getDbFunc: getDbFunc,
	}
}

func (s *LoanPolicyTemplateRepository) GetById(ctx context.Context, id int64) (entity.LoanPolicyTemplate, error) {
	var dest model.LoanPolicyTemplate
	if err := table.LoanPolicyTemplate.
		SELECT(table.LoanPolicyTemplate.AllColumns).
		WHERE(table.LoanPolicyTemplate.ID.EQ(postgres.Int64(id))).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		return entity.LoanPolicyTemplate{}, fmt.Errorf("LoanPolicyTemplateRepository GetById %w", err)
	}
	return MapLoanPolicyTemplateDbToEntity(dest), nil
}

func (s *LoanPolicyTemplateRepository) GetByIds(ctx context.Context, ids []int64) ([]entity.LoanPolicyTemplate, error) {
	var dest []model.LoanPolicyTemplate
	if err := table.LoanPolicyTemplate.
		SELECT(table.LoanPolicyTemplate.AllColumns).
		WHERE(table.LoanPolicyTemplate.ID.IN(querymod.In(ids)...)).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		return []entity.LoanPolicyTemplate{}, fmt.Errorf("LoanPolicyTemplateRepository GetByIds %w", err)
	}
	return MapLoanPolicyTemplatesDbToEntities(dest), nil
}

func (s *LoanPolicyTemplateRepository) Delete(ctx context.Context, id int64) error {
	if _, err := table.LoanPolicyTemplate.DELETE().
		WHERE(table.LoanPolicyTemplate.ID.EQ(postgres.Int64(id))).ExecContext(
		ctx, s.getDbFunc(ctx),
	); err != nil {
		return fmt.Errorf("LoanPolicyTemplateRepository Delete %w", err)
	}
	return nil
}

func (s *LoanPolicyTemplateRepository) GetAll(ctx context.Context) ([]entity.LoanPolicyTemplate, error) {
	dest := make([]model.LoanPolicyTemplate, 0)
	if err := table.LoanPolicyTemplate.SELECT(table.LoanPolicyTemplate.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []entity.LoanPolicyTemplate{}, nil
		}
		return []entity.LoanPolicyTemplate{}, fmt.Errorf("LoanPolicyTemplateRepository GetAll %w", err)
	}
	return MapLoanPolicyTemplatesDbToEntities(dest), nil
}

func (s *LoanPolicyTemplateRepository) Create(ctx context.Context, template entity.LoanPolicyTemplate) (entity.LoanPolicyTemplate, error) {
	createModel := MapLoanPolicyTemplateEntityToDb(template)
	created := model.LoanPolicyTemplate{}
	if err := table.LoanPolicyTemplate.INSERT(table.LoanPolicyTemplate.MutableColumns).
		MODEL(createModel).
		RETURNING(table.LoanPolicyTemplate.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &created); err != nil {
		return entity.LoanPolicyTemplate{}, fmt.Errorf("LoanPolicyTemplateRepository Create %w", err)
	}
	return MapLoanPolicyTemplateDbToEntity(created), nil
}

func (s *LoanPolicyTemplateRepository) Update(ctx context.Context, template entity.LoanPolicyTemplate) (entity.LoanPolicyTemplate, error) {
	updateModel := MapLoanPolicyTemplateEntityToDb(template)
	updated := model.LoanPolicyTemplate{}
	if err := table.LoanPolicyTemplate.UPDATE(table.LoanPolicyTemplate.MutableColumns).
		MODEL(updateModel).
		WHERE(table.LoanPolicyTemplate.ID.EQ(postgres.Int64(template.Id))).
		RETURNING(table.LoanPolicyTemplate.AllColumns).
		QueryContext(ctx, s.getDbFunc(ctx), &updated); err != nil {
		return entity.LoanPolicyTemplate{}, fmt.Errorf("LoanPolicyTemplateRepository Update %w", err)
	}
	return MapLoanPolicyTemplateDbToEntity(updated), nil
}
