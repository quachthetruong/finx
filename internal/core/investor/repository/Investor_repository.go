package repository

import (
	"context"

	"financing-offer/internal/core/entity"
)

type InvestorPersistenceRepository interface {
	BulkCreate(ctx context.Context, investors []entity.Investor) error
	CreateIfNotExist(ctx context.Context, investor entity.Investor) error
	Update(ctx context.Context, investor entity.Investor) (entity.Investor, error)
	GetAllUniqueInvestorIdsFromRequests(ctx context.Context) ([]string, error)
	GetAllInvestorIdsForMigration(ctx context.Context) ([]string, error)
}
