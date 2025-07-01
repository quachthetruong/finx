package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database"
	"financing-offer/pkg/dbtest"
)

func TestSuggestedOfferRepository_Create(t *testing.T) {
	t.Parallel()
	db, mock, err := dbtest.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	repo := NewSuggestedOfferRepository(
		func(ctx context.Context) database.DB {
			return db
		},
	)
	// Test data
	suggestedOffer := entity.SuggestedOffer{
		Id:        1,
		ConfigId:  1,
		AccountNo: "0001000115",
		Symbols: []string{
			"ACB",
			"VCB",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Test case
	t.Run("create suggested offer success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"suggested_offer.id",
			"suggested_offer.config_id",
			"suggested_offer.account_no",
			"suggested_offer.symbols",
			"suggested_offer.created_at",
			"suggested_offer.updated_at",
		})
		rows.AddRow(
			suggestedOffer.Id,
			suggestedOffer.ConfigId,
			suggestedOffer.AccountNo,
			"[\"ACB\", \"VCB\"]",
			suggestedOffer.CreatedAt,
			suggestedOffer.UpdatedAt,
		)
		mock.ExpectQuery("INSERT").WillReturnRows(rows)
		res, err := repo.Create(context.Background(), suggestedOffer)
		assert.Nil(t, err)
		assert.Equal(t, suggestedOffer, res)
	})
	t.Run("create suggested offer error", func(t *testing.T) {
		mock.ExpectQuery("INSERT").WillReturnError(assert.AnError)
		res, err := repo.Create(context.Background(), suggestedOffer)
		assert.Empty(t, res)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
