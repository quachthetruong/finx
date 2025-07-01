package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
)

func TestLoanTemplateRepository(t *testing.T) {
	t.Parallel()
	db, mock, _ := dbtest.New()
	repo := NewLoanPolicyTemplateRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)

	t.Run("Create", func(t *testing.T) {
		template := entity.LoanPolicyTemplate{
			Id:                       1,
			CreatedAt:                time.Now(),
			UpdatedAt:                time.Now(),
			Name:                     "test",
			InterestRate:             decimal.NewFromInt(3),
			InterestBasis:            1,
			Term:                     1,
			PoolIdRef:                1,
			OverdueInterest:          decimal.NewFromInt(3),
			AllowExtendLoanTerm:      true,
			AllowEarlyPayment:        true,
			PreferentialPeriod:       1,
			PreferentialInterestRate: decimal.NewFromInt(3),
		}
		mock.ExpectQuery("INSERT").WillReturnRows(
			mock.NewRows(
				[]string{
					"loan_policy_template.id",
					"loan_policy_template.created_at",
					"loan_policy_template.updated_at",
					"loan_policy_template.updated_by",
					"loan_policy_template.name",
					"loan_policy_template.interest_rate",
					"loan_policy_template.interest_basis",
					"loan_policy_template.term",
					"loan_policy_template.pool_id_ref",
					"loan_policy_template.overdue_interest",
					"loan_policy_template.allow_extend_loan_term",
					"loan_policy_template.allow_early_payment",
					"loan_policy_template.preferential_period",
					"loan_policy_template.preferential_interest_rate",
				}).AddRow(template.Id, template.CreatedAt, template.UpdatedAt, template.UpdatedBy,
				template.Name, template.InterestRate, template.InterestBasis, template.Term,
				template.PoolIdRef, template.OverdueInterest, template.AllowExtendLoanTerm, template.AllowEarlyPayment,
				template.PreferentialPeriod, template.PreferentialInterestRate))
		result, err := repo.Create(context.Background(), template)
		assert.Nil(t, err)
		assert.Equal(t, template.Id, result.Id)
		assert.Equal(t, template.Name, result.Name)
	})

	t.Run("GetAll", func(t *testing.T) {
		policies := []entity.LoanPolicyTemplate{
			{
				Id:                       1,
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
				UpdatedBy:                "admin",
				Name:                     "test",
				InterestRate:             decimal.NewFromInt(3),
				InterestBasis:            1,
				Term:                     1,
				PoolIdRef:                1,
				OverdueInterest:          decimal.NewFromInt(3),
				AllowExtendLoanTerm:      true,
				AllowEarlyPayment:        true,
				PreferentialPeriod:       1,
				PreferentialInterestRate: decimal.NewFromInt(3),
			},
			{
				Id:                       2,
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
				UpdatedBy:                "admin",
				Name:                     "test2",
				InterestRate:             decimal.NewFromInt(32),
				InterestBasis:            12,
				Term:                     12,
				PoolIdRef:                12,
				OverdueInterest:          decimal.NewFromInt(32),
				AllowExtendLoanTerm:      true,
				AllowEarlyPayment:        true,
				PreferentialPeriod:       12,
				PreferentialInterestRate: decimal.NewFromInt(32),
			},
		}
		mock.ExpectQuery("SELECT").WillReturnRows(
			mock.NewRows(
				[]string{
					"loan_policy_template.id",
					"loan_policy_template.created_at",
					"loan_policy_template.updated_at",
					"loan_policy_template.updated_by",
					"loan_policy_template.name",
					"loan_policy_template.interest_rate",
					"loan_policy_template.interest_basis",
					"loan_policy_template.term",
					"loan_policy_template.pool_id_ref",
					"loan_policy_template.overdue_interest",
					"loan_policy_template.allow_extend_loan_term",
					"loan_policy_template.allow_early_payment",
					"loan_policy_template.preferential_period",
					"loan_policy_template.preferential_interest_rate",
				}).AddRow(policies[0].Id, policies[0].CreatedAt, policies[0].UpdatedAt, policies[0].UpdatedBy,
				policies[0].Name, policies[0].InterestRate, policies[0].InterestBasis, policies[0].Term,
				policies[0].PoolIdRef, policies[0].OverdueInterest, policies[0].AllowExtendLoanTerm, policies[0].AllowEarlyPayment,
				policies[0].PreferentialPeriod, policies[0].PreferentialInterestRate).
				AddRow(policies[1].Id, policies[1].CreatedAt, policies[1].UpdatedAt, policies[1].UpdatedBy,
					policies[1].Name, policies[1].InterestRate, policies[1].InterestBasis, policies[1].Term,
					policies[1].PoolIdRef, policies[1].OverdueInterest, policies[1].AllowExtendLoanTerm, policies[1].AllowEarlyPayment,
					policies[1].PreferentialPeriod, policies[1].PreferentialInterestRate))
		result, err := repo.GetAll(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, len(result), 2)
		assert.Equal(t, result[0].Id, int64(1))
		assert.Equal(t, result[1].Id, int64(2))
	})

	t.Run("Update", func(t *testing.T) {
		template := entity.LoanPolicyTemplate{
			Id:                       1,
			CreatedAt:                time.Now(),
			UpdatedAt:                time.Now(),
			Name:                     "test",
			InterestRate:             decimal.NewFromInt(3),
			InterestBasis:            1,
			Term:                     1,
			PoolIdRef:                1,
			OverdueInterest:          decimal.NewFromInt(3),
			AllowExtendLoanTerm:      true,
			AllowEarlyPayment:        true,
			PreferentialPeriod:       1,
			PreferentialInterestRate: decimal.NewFromInt(3),
		}
		mock.ExpectQuery("UPDATE").WillReturnRows(
			mock.NewRows(
				[]string{
					"loan_policy_template.id",
					"loan_policy_template.created_at",
					"loan_policy_template.updated_at",
					"loan_policy_template.updated_by",
					"loan_policy_template.name",
					"loan_policy_template.interest_rate",
					"loan_policy_template.interest_basis",
					"loan_policy_template.term",
					"loan_policy_template.pool_id_ref",
					"loan_policy_template.overdue_interest",
					"loan_policy_template.allow_extend_loan_term",
					"loan_policy_template.allow_early_payment",
					"loan_policy_template.preferential_period",
					"loan_policy_template.preferential_interest_rate",
				}).AddRow(template.Id, template.CreatedAt, template.UpdatedAt, template.UpdatedBy,
				template.Name, template.InterestRate, template.InterestBasis, template.Term,
				template.PoolIdRef, template.OverdueInterest, template.AllowExtendLoanTerm, template.AllowEarlyPayment,
				template.PreferentialPeriod, template.PreferentialInterestRate))
		result, err := repo.Update(context.Background(), template)
		assert.Nil(t, err)
		assert.Equal(t, template.Id, result.Id)
		assert.Equal(t, template.Name, result.Name)
	})
}
